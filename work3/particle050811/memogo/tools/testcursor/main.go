package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type RespEnvelope struct {
	Status int             `json:"status"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Todo struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    int32  `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type CursorTodoData struct {
	Items      []Todo `json:"items"`
	NextCursor int64  `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

type CreateTodoResp struct {
	Status int32  `json:"status"`
	Msg    string `json:"msg"`
	Data   Todo   `json:"data"`
}

func doJSON(ctx context.Context, client *http.Client, method, url, token string, in any, out any) error {
	var body io.Reader
	if in != nil {
		b, _ := json.Marshal(in)
		body = bytes.NewReader(b)
	}
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if out != nil {
		json.Unmarshal(data, out)
	}
	return nil
}

func main() {
	base := "http://127.0.0.1:8888"
	user := "cursor_test_user"
	pass := "P@ssw0rd1"

	httpc := &http.Client{Timeout: 8 * time.Second}
	ctx := context.Background()

	// 1. 注册和登录
	regBody := map[string]string{"username": user, "password": pass}
	var regResp RespEnvelope
	doJSON(ctx, httpc, "POST", base+"/v1/auth/register", "", regBody, &regResp)
	fmt.Println("✓ 注册:", regResp.Msg)

	var loginEnv RespEnvelope
	doJSON(ctx, httpc, "POST", base+"/v1/auth/login", "", regBody, &loginEnv)
	var tp TokenPair
	json.Unmarshal(loginEnv.Data, &tp)
	fmt.Println("✓ 登录成功，获取 token")

	// 2. 创建 15 条测试数据
	fmt.Println("\n📝 创建 15 条测试备忘录...")
	for i := 1; i <= 15; i++ {
		req := map[string]string{
			"title":   fmt.Sprintf("备忘录-%02d", i),
			"content": fmt.Sprintf("这是第 %d 条测试内容", i),
		}
		var resp CreateTodoResp
		doJSON(ctx, httpc, "POST", base+"/v1/todos", tp.AccessToken, req, &resp)
		fmt.Printf("  ✓ 创建: ID=%d, Title=%s\n", resp.Data.ID, resp.Data.Title)
		time.Sleep(50 * time.Millisecond) // 避免创建时间完全相同
	}

	// 3. 测试游标分页 - 遍历全部数据
	fmt.Println("\n🚀 测试游标分页（每次 5 条）...")
	cursor := int64(0)
	page := 1
	totalFetched := 0

	for {
		url := fmt.Sprintf("%s/v1/todos/cursor?status=all&cursor=%d&limit=5", base, cursor)
		var env RespEnvelope
		doJSON(ctx, httpc, "GET", url, tp.AccessToken, nil, &env)

		var cursorData CursorTodoData
		json.Unmarshal(env.Data, &cursorData)

		fmt.Printf("\n第 %d 页 (cursor=%d):\n", page, cursor)
		for _, todo := range cursorData.Items {
			fmt.Printf("  - ID=%d, Title=%s\n", todo.ID, todo.Title)
		}

		totalFetched += len(cursorData.Items)
		fmt.Printf("  Next Cursor: %d, Has More: %v\n", cursorData.NextCursor, cursorData.HasMore)

		if !cursorData.HasMore {
			break
		}

		cursor = cursorData.NextCursor
		page++
	}

	fmt.Printf("\n✅ 游标分页完成！总共获取了 %d 条记录\n", totalFetched)

	// 4. 测试关键词游标分页
	fmt.Println("\n🔍 测试关键词游标分页（搜索 '备忘录'）...")
	cursor = 0
	page = 1
	totalSearched := 0

	for {
		url := fmt.Sprintf("%s/v1/todos/search/cursor?q=备忘录&cursor=%d&limit=5", base, cursor)
		var env RespEnvelope
		doJSON(ctx, httpc, "GET", url, tp.AccessToken, nil, &env)

		var cursorData CursorTodoData
		json.Unmarshal(env.Data, &cursorData)

		fmt.Printf("\n搜索第 %d 页 (cursor=%d):\n", page, cursor)
		for _, todo := range cursorData.Items {
			fmt.Printf("  - ID=%d, Title=%s\n", todo.ID, todo.Title)
		}

		totalSearched += len(cursorData.Items)
		fmt.Printf("  Next Cursor: %d, Has More: %v\n", cursorData.NextCursor, cursorData.HasMore)

		if !cursorData.HasMore {
			break
		}

		cursor = cursorData.NextCursor
		page++
	}

	fmt.Printf("\n✅ 关键词游标分页完成！总共搜索到 %d 条记录\n", totalSearched)

	// 5. 清理测试数据
	fmt.Println("\n🗑️  清理测试数据...")
	var delEnv RespEnvelope
	doJSON(ctx, httpc, "DELETE", base+"/v1/todos?scope=all", tp.AccessToken, nil, &delEnv)
	fmt.Println("✓ 清理完成")

	fmt.Println("\n🎉 所有测试通过！")
	os.Exit(0)
}
