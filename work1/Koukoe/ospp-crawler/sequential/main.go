package main

import (
	"fmt"
	"net/http"
	"os"
	"ospp-crawler/common"
	"time"
)

func crawl(client *http.Client) ([]common.Project, error) {
	var projects []common.Project
	for page := 1; ; page++ {
		// if page > 2 {
		// 	break
		// }

		rows, total, err := common.FetchPage(client, page)
		if err != nil {
			return nil, err
		}

		if len(rows) == 0 {
			break
		}

		projects = append(projects, rows...)
		fmt.Printf("第 %d 页：本页 %d 条，累计 %d/%d\n", page, len(rows), len(projects), total)

		if len(projects) >= total {
			break
		}
	}
	return projects, nil
}

func main() {
	start := time.Now()
	client := common.NewClient()

	projects, err := crawl(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, "爬取失败:", err)
		os.Exit(1)
	}

	n, err := common.Save(projects)
	if err != nil {
		fmt.Fprintln(os.Stderr, "保存失败:", err)
		os.Exit(1)
	}

	fmt.Printf("\n共 %d 个项目, 已保存到 %s, 用时 %v\n", n, common.OutFile, time.Since(start))

	downloadApplications(client, projects)
	fmt.Printf("\n完成, 总用时 %v\n", time.Since(start))
}

func downloadApplications(client *http.Client, projects []common.Project) {
	var count int

	for i, proj := range projects {
		data, err := common.FetchApplication(client, proj.ProId)
		switch {
		case err != nil:
			fmt.Printf("[%d/%d] %s 获取失败: %v\n", i+1, len(projects), proj.ProgramCode, err)
			continue
		case data == nil:
			fmt.Printf("[%d/%d] %s 没有申请书\n", i+1, len(projects), proj.ProgramCode)
			continue
		}

		path, err := common.SaveApplication(proj, data)
		if err != nil {
			fmt.Printf("[%d/%d] %s 保存失败: %v\n", i+1, len(projects), proj.ProgramCode, err)
			continue
		}
		count++
		fmt.Printf("[%d/%d] %s 保存到 %s (%d KB)\n", i+1, len(projects), proj.ProgramCode, path, len(data)/1024)
	}

	fmt.Printf("\n%d 个申请书下载完成\n", count) //	好多
}
