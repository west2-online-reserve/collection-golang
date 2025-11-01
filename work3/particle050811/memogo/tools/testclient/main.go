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

func doJSON(ctx context.Context, client *http.Client, method, url, token string, in any, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
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
	flag.Parse()

	httpc := &http.Client{Timeout: 8 * time.Second}
	ctx := context.Background()

	// 0) wait server ready (retry ping for up to 10s)
	if err := waitServerReady(ctx, httpc, *base, 10*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, "Server not ready:", err)
		os.Exit(1)
	}

	// 1) ping (final confirmation)
	pingURL := *base + "/ping"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, pingURL, nil)
	if resp, err := httpc.Do(req); err == nil {
		resp.Body.Close()
		fmt.Println("Ping ok:", pingURL)
	} else {
		fmt.Fprintln(os.Stderr, "Ping failed:", err)
	}

	// 2) register (ignore already exists)
	regURL := *base + "/v1/auth/register"
	regBody := map[string]string{"username": *user, "password": *pass}
	var regResp RespEnvelope
	_ = doJSON(ctx, httpc, http.MethodPost, regURL, "", regBody, &regResp)
	fmt.Println("Register status:", regResp.Status, regResp.Msg)

	// 3) login
	loginURL := *base + "/v1/auth/login"
	var loginEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodPost, loginURL, "", regBody, &loginEnv); err != nil {
		fmt.Fprintln(os.Stderr, "Login failed:", err)
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
		fmt.Fprintln(os.Stderr, "No access_token in login response")
		os.Exit(1)
	}
	fmt.Println("Login ok, token acquired")

	// 4) create 3 todos
	createURL := *base + "/v1/todos"
	createdIDs := make([]int64, 0, 3)
	for i := 1; i <= 3; i++ {
		title := fmt.Sprintf("测试任务-%d", i)
		cReq := CreateTodoReq{Title: title, Content: *keyword}
		var cEnv CreateTodoResp
		if err := doJSON(ctx, httpc, http.MethodPost, createURL, tp.AccessToken, cReq, &cEnv); err != nil {
			fmt.Fprintln(os.Stderr, "Create todo failed:", err)
			os.Exit(1)
		}
		createdIDs = append(createdIDs, cEnv.Data.ID)
		fmt.Println("Created:", title)
	}

	// 5) list todos (all) page=1 size=2
	listURL := *base + "/v1/todos?page=1&page_size=2"
	var listEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodGet, listURL, tp.AccessToken, nil, &listEnv); err != nil {
		fmt.Fprintln(os.Stderr, "List failed:", err)
		os.Exit(1)
	}
	var listData ItemsTodoData
	_ = json.Unmarshal(listEnv.Data, &listData)
	fmt.Printf("List: total=%d, pageItems=%d\n", listData.Total, len(listData.Items))

	// 6) search by keyword
	searchURL := *base + "/v1/todos/search?q=" + *keyword + "&page=1&page_size=10"
	var searchEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodGet, searchURL, tp.AccessToken, nil, &searchEnv); err != nil {
		fmt.Fprintln(os.Stderr, "Search failed:", err)
		os.Exit(1)
	}
	var searchData ItemsTodoData
	_ = json.Unmarshal(searchEnv.Data, &searchData)
	fmt.Printf("Search '%s': total=%d, pageItems=%d\n", *keyword, searchData.Total, len(searchData.Items))

	// 7) update one to DONE if any (prefer use the last created id)
	if len(createdIDs) > 0 {
		id := createdIDs[len(createdIDs)-1]
		updURL := fmt.Sprintf("%s/v1/todos/%d/status", *base, id)
		body := map[string]int{"status": 1}
		fmt.Println("Updating one via:", updURL)

		// ▶ 只加这一行：把 HTTP 请求关键信息一次性打印出来
		//fmt.Printf("HTTP PATCH %s\nAuthorization: Bearer %s\nContent-Type: application/json\n\n%v\n",
		//	updURL, tp.AccessToken, body)

		var updEnv RespEnvelope
		if err := doJSON(ctx, httpc, http.MethodPatch, updURL, tp.AccessToken, body, &updEnv); err != nil {
			fmt.Fprintln(os.Stderr, "Update one failed:", err)
			// 继续后续批量步骤，避免单点失败中断整个流程
		} else {
			fmt.Println("Updated one to DONE, id=", id)
		}

	}

	// 8) batch TODO->DONE
	var batchEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodPatch, *base+"/v1/todos/status?from=0&to=1", tp.AccessToken, nil, &batchEnv); err != nil {
		fmt.Fprintln(os.Stderr, "Batch TODO->DONE failed:", err)
		os.Exit(1)
	}
	fmt.Println("Batch TODO->DONE ok")

	// 9) batch DONE->TODO
	if err := doJSON(ctx, httpc, http.MethodPatch, *base+"/v1/todos/status?from=1&to=0", tp.AccessToken, nil, &batchEnv); err != nil {
		fmt.Fprintln(os.Stderr, "Batch DONE->TODO failed:", err)
		os.Exit(1)
	}
	fmt.Println("Batch DONE->TODO ok")

	// 10) delete by scope
	var delEnv RespEnvelope
	if err := doJSON(ctx, httpc, http.MethodDelete, *base+"/v1/todos?scope=done", tp.AccessToken, nil, &delEnv); err != nil {
		fmt.Fprintln(os.Stderr, "Delete done failed:", err)
		os.Exit(1)
	}
	if err := doJSON(ctx, httpc, http.MethodDelete, *base+"/v1/todos?scope=todo", tp.AccessToken, nil, &delEnv); err != nil {
		fmt.Fprintln(os.Stderr, "Delete todo failed:", err)
		os.Exit(1)
	}
	fmt.Println("Delete by scope done+todo ok")

	// 11) 测试单条删除 - 重新创建一条待办事项并删除
	fmt.Println("\n--- 测试单条删除 ---")
	var singleTodoResp CreateTodoResp
	singleReq := CreateTodoReq{Title: "测试单条删除", Content: "这条待办事项将被删除"}
	if err := doJSON(ctx, httpc, http.MethodPost, createURL, tp.AccessToken, singleReq, &singleTodoResp); err != nil {
		fmt.Fprintln(os.Stderr, "Create single todo for delete test failed:", err)
		os.Exit(1)
	}
	fmt.Println("Created single todo for delete test, ID:", singleTodoResp.Data.ID)

	// 删除单条
	deleteOneURL := fmt.Sprintf("%s/v1/todos/%d", *base, singleTodoResp.Data.ID)
	fmt.Println("Deleting single todo via:", deleteOneURL)
	if err := doJSON(ctx, httpc, http.MethodDelete, deleteOneURL, tp.AccessToken, nil, &delEnv); err != nil {
		fmt.Fprintln(os.Stderr, "Delete single todo failed:", err)
		os.Exit(1)
	}
	fmt.Println("Delete single todo ok")

	fmt.Println("Scenario completed successfully.")
}

// waitServerReady keeps probing GET {base}/ping until success or timeout.
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
		// exponential backoff up to 1s
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
