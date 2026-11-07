package main

import (
	"fmt"
	"os"
	"ospp-crawler/common"
	"time"
)

func crawl() ([]common.Project, error) {
	client := common.NewClient()

	var projects []common.Project
	for page := 1; ; page++ {
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

	projects, err := crawl()
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
}
