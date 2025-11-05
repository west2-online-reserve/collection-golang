package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	View      int32  `json:"view"`
	Status    int32  `json:"status"` // 0 TODO, 1 DONE
	CreatedAt int64  `json:"created_at"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	DueTime   int64  `json:"due_time"`
}

type ItemsTodoData struct {
	Items []Todo `json:"items"`
	Total int64  `json:"total"`
}

type CreateTodoReq struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	StartTime *int64 `json:"start_time,omitempty"`
	DueTime   *int64 `json:"due_time,omitempty"`
}

type CreateTodoResp struct {
	Status int32  `json:"status"`
	Msg    string `json:"msg"`
	Data   Todo   `json:"data"`
}

var verbose bool // 全局变量，控制是否输出详细日志

func doJSON(ctx context.Context, client *http.Client, method, url, token string, in any, out any) error {
	var body io.Reader
	var reqBody []byte
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		reqBody = b
		body = bytes.NewReader(b)
	}

	// 详细日志：请求信息
	if verbose {
		fmt.Printf("\n==> %s %s\n", method, url)
		if token != "" {
			fmt.Printf("    Authorization: Bearer %s...\n", token[:min(20, len(token))])
		}
		if reqBody != nil {
			fmt.Printf("    Request Body: %s\n", string(reqBody))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
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
	data, _ := io.ReadAll(resp.Body)

	// 详细日志：响应信息
	if verbose {
		fmt.Printf("<== HTTP %d\n", resp.StatusCode)
		fmt.Printf("    Response Body: %s\n", strings.TrimSpace(string(data)))
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if out != nil {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("unmarshal: %w (body=%s)", err, string(data))
		}
	}
	return nil
}

func main() {
	base := flag.String("base", "http://127.0.0.1:8888", "API base URL")
	user := flag.String("user", "tester", "username")
	pass := flag.String("pass", "P@ssw0rd1", "password")
	keyword := flag.String("q", "keywordX", "search keyword")
	verboseFlag := flag.Bool("v", false, "verbose mode: print detailed HTTP requests/responses")
	flag.Parse()

	verbose = *verboseFlag

	httpc := &http.Client{Timeout: 8 * time.Second}
	ctx := context.Background()

	// 0) 等待服务就绪（重试 ping，最多 10 秒）
	if err := waitServerReady(ctx, httpc, *base, 10*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, "服务未就绪:", err)
		os.Exit(1)
	}

	// 1) ping（最终确认）
	pingURL := *base + "/ping"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, pingURL, nil)
	if resp, err := httpc.Do(req); err == nil {
		resp.Body.Close()
		fmt.Println("Ping 成功:", pingURL)
	} else {
		fmt.Fprintln(os.Stderr, "Ping 失败:", err)
	}

	// 2) 注册（如果用户已存在则忽略）
	regURL := *base + "/v1/auth/register"
	regBody := map[string]string{"username": *user, "password": *pass}
	var regResp RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodPost, regURL, "", regBody, &regResp); err != nil {
		// 用户已存在是正常情况，继续登录即可
		if strings.Contains(err.Error(), "HTTP 400") && strings.Contains(err.Error(), "already exists") {
			fmt.Printf("注册: 用户 '%s' 已存在，跳过注册\n", *user)
		} else {
			fmt.Printf("注册失败: %v\n", err)
		}
	} else {
		fmt.Printf("注册成功: %s (status=%d)\n", regResp.Msg, regResp.Status)
	}

	// 3) 登录
	loginURL := *base + "/v1/auth/login"
	var loginEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodPost, loginURL, "", regBody, &loginEnv); err != nil {
		fmt.Fprintln(os.Stderr, "登录失败:", err)
		os.Exit(1)
	}
	var tp TokenPair
	if len(loginEnv.Data) > 0 {
		_ = json.Unmarshal(loginEnv.Data, &tp)
	}
	if tp.AccessToken == "" {
		// 有些实现直接把 token 放在顶层 data 字段里，这里兜底再解析一遍
		_ = json.Unmarshal(loginEnv.Data, &tp)
	}
	if tp.AccessToken == "" {
		fmt.Fprintln(os.Stderr, "登录响应中没有 access_token")
		os.Exit(1)
	}
	fmt.Println("登录成功，已获取 token")

	// 4) 创建 3 条待办事项
	fmt.Println("\n创建待办测试:")
	createURL := *base + "/v1/todos"
	createdIDs := make([]int64, 0, 3)
	for i := 1; i <= 3; i++ {
		title := fmt.Sprintf("测试任务-%d", i)
		cReq := CreateTodoReq{Title: title, Content: *keyword}
		var cEnv CreateTodoResp
		if err := doJSON(ctx, httpc, http.MethodPost, createURL, tp.AccessToken, cReq, &cEnv); err != nil {
			fmt.Fprintln(os.Stderr, "  ✗ 创建待办失败:", err)
			os.Exit(1)
		}
		createdIDs = append(createdIDs, cEnv.Data.ID)
		fmt.Printf("  ✓ 已创建: '%s' (ID=%d)\n", title, cEnv.Data.ID)
	}

	// 5) 列出待办事项（全部）page=1 size=2
	fmt.Println("\n查询待办测试:")
	listURL := *base + "/v1/todos?page=1&page_size=2"
	var listEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodGet, listURL, tp.AccessToken, nil, &listEnv); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 列表查询失败:", err)
		os.Exit(1)
	}
	var listData ItemsTodoData
	_ = json.Unmarshal(listEnv.Data, &listData)
	fmt.Printf("  ✓ 列表查询 (总数=%d, 当前页=%d 条)\n", listData.Total, len(listData.Items))

	// 6) 按关键词搜索
	searchURL := *base + "/v1/todos/search?q=" + *keyword + "&page=1&page_size=10"
	var searchEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodGet, searchURL, tp.AccessToken, nil, &searchEnv); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 搜索失败:", err)
		os.Exit(1)
	}
	var searchData ItemsTodoData
	_ = json.Unmarshal(searchEnv.Data, &searchData)
	fmt.Printf("  ✓ 关键词搜索 '%s' (总数=%d, 当前页=%d 条)\n", *keyword, searchData.Total, len(searchData.Items))

	// 7) 更新一条为 DONE（优先使用最后创建的 id）
	fmt.Println("\n更新单条状态测试:")
	if len(createdIDs) > 0 {
		id := createdIDs[len(createdIDs)-1]
		updURL := fmt.Sprintf("%s/v1/todos/%d/status", *base, id)
		body := map[string]int{"status": 1}

		var updEnv RespEnvelope
		if err := doJSON(ctx, httpc, http.MethodPatch, updURL, tp.AccessToken, body, &updEnv); err != nil {
			fmt.Fprintln(os.Stderr, "  ✗ 更新单条失败:", err)
			// 继续后续批量步骤，避免单点失败中断整个流程
		} else {
			var affected int
			_ = json.Unmarshal(updEnv.Data, &affected)
			fmt.Printf("  ✓ TODO -> DONE (ID=%d, 受影响: %d 条)\n", id, affected)
		}
	}

	// 8) 批量 TODO->DONE
	fmt.Println("\n批量更新状态测试:")
	var batchEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodPatch, *base+"/v1/todos/status?from=0&to=1", tp.AccessToken, nil, &batchEnv); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 批量 TODO->DONE 失败:", err)
		os.Exit(1)
	}
	var affected int
	_ = json.Unmarshal(batchEnv.Data, &affected)
	fmt.Printf("  ✓ 批量 TODO->DONE 成功 (受影响: %d 条)\n", affected)

	// 9) 批量 DONE->TODO
	if err := doJSON(ctx, httpc, http.MethodPatch, *base+"/v1/todos/status?from=1&to=0", tp.AccessToken, nil, &batchEnv); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 批量 DONE->TODO 失败:", err)
		os.Exit(1)
	}
	_ = json.Unmarshal(batchEnv.Data, &affected)
	fmt.Printf("  ✓ 批量 DONE->TODO 成功 (受影响: %d 条)\n", affected)

	// 10) 按范围删除
	fmt.Println("\n按范围删除测试:")
	var delEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodDelete, *base+"/v1/todos?scope=done", tp.AccessToken, nil, &delEnv); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 删除已完成项失败:", err)
		os.Exit(1)
	}
	var deletedDone int
	_ = json.Unmarshal(delEnv.Data, &deletedDone)
	fmt.Printf("  ✓ 删除已完成项 (删除: %d 条)\n", deletedDone)

	if err := doJSON(ctx, httpc, http.MethodDelete, *base+"/v1/todos?scope=todo", tp.AccessToken, nil, &delEnv); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 删除待办项失败:", err)
		os.Exit(1)
	}
	var deletedTodo int
	_ = json.Unmarshal(delEnv.Data, &deletedTodo)
	fmt.Printf("  ✓ 删除待办项 (删除: %d 条)\n", deletedTodo)

	// 11) 测试单条删除 - 重新创建一条待办事项并删除
	fmt.Println("\n单条删除测试:")
	var singleTodoResp CreateTodoResp
	singleReq := CreateTodoReq{Title: "测试单条删除", Content: "这条待办事项将被删除"}
	if err := doJSON(ctx, httpc, http.MethodPost, createURL, tp.AccessToken, singleReq, &singleTodoResp); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 创建单条待办失败:", err)
		os.Exit(1)
	}
	fmt.Printf("  ➤ 已创建待办 (ID=%d, 标题='%s')\n", singleTodoResp.Data.ID, singleTodoResp.Data.Title)

	// 删除单条
	deleteOneURL := fmt.Sprintf("%s/v1/todos/%d", *base, singleTodoResp.Data.ID)
	if err := doJSON(ctx, httpc, http.MethodDelete, deleteOneURL, tp.AccessToken, nil, &delEnv); err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ 删除单条待办失败:", err)
		os.Exit(1)
	}
	var deletedSingle int
	_ = json.Unmarshal(delEnv.Data, &deletedSingle)
	fmt.Printf("  ✓ 删除单条待办成功 (ID=%d, 受影响: %d 条)\n", singleTodoResp.Data.ID, deletedSingle)

	fmt.Println("\n✅ 所有测试场景执行成功")
}

// waitServerReady 持续探测 GET {base}/ping 直到成功或超时
func waitServerReady(ctx context.Context, c *http.Client, base string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	url := base + "/ping"
	for attempt := 0; ; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if resp, err := c.Do(req); err == nil {
			resp.Body.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for %s", url)
		}
		// 指数退避，最多 1 秒
		sleep := 100 * time.Millisecond * (1 << min(attempt, 3))
		time.Sleep(sleep)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
