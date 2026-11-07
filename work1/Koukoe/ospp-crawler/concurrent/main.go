package main

import (
	"fmt"
	"net/http"
	"os"
	"ospp-crawler/common"
	"sync"
	"time"
)

const workers = 10 //  加速比差不多是3.8

func crawl(client *http.Client) ([]common.Project, error) {
	first, total, err := common.FetchPage(client, 1)
	if err != nil {
		return nil, err
	}

	pageCount := (total + common.PageSize - 1) / common.PageSize
	if pageCount == 0 {
		return first, nil
	}

	pages := make([][]common.Project, pageCount)
	errs := make([]error, pageCount)
	pages[0] = first

	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range jobs {
				rows, _, err := common.FetchPage(client, page)
				pages[page-1] = rows
				errs[page-1] = err
			}
		}()
	}

	for page := 2; page <= pageCount; page++ {
		jobs <- page
	}
	close(jobs)
	wg.Wait()

	projects := make([]common.Project, 0, total)
	for i, rows := range pages {
		if errs[i] != nil {
			return nil, errs[i]
		}
		projects = append(projects, rows...)
		fmt.Printf("第 %d 页：本页 %d 条，累计 %d/%d\n", i+1, len(rows), len(projects), total)
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
}
