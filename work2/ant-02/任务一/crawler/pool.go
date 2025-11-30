package crawler

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"west2/db"
	"west2/model"
	"west2/util"

	"golang.org/x/net/html"
)

// Pool 协程池结构
type Pool struct {
	taskChan    chan *Task         // 任务通道
	workerCount int32              // 当前worker数量
	maxWorkers  int                // 最大worker数量
	wg          sync.WaitGroup     // 等待组
	ctx         context.Context    // 上下文
	cancel      context.CancelFunc // 取消函数
	isStopped   int32              // 停止标志
	stats       *PoolStats         // 统计信息
}

// Task 任务结构
type Task struct {
	URL     string
	PageNum int
	Type    TaskType
	Data    interface{}
}

// TaskType 任务类型
type TaskType int

const (
	TaskTypePage TaskType = iota // 页面任务
	TaskTypeItem                 // 列表项任务
)

// PoolStats 池统计
type PoolStats struct {
	TotalTasks     int64
	ProcessedTasks int64
	FailedTasks    int64
	mu             sync.RWMutex
}

// NewPool 创建新的协程池
func NewPool(maxWorkers int, queueSize int) *Pool {
	if maxWorkers <= 0 {
		maxWorkers = runtime.NumCPU() * 2
	}
	if queueSize <= 0 {
		queueSize = 1000
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		taskChan:   make(chan *Task, queueSize),
		maxWorkers: maxWorkers,
		ctx:        ctx,
		cancel:     cancel,
		stats:      &PoolStats{},
	}

	// 启动初始worker
	pool.startWorkers(runtime.NumCPU())

	// 启动监控
	go pool.monitor()

	return pool
}

// startWorkers 启动指定数量的worker
func (p *Pool) startWorkers(count int) {
	for i := 0; i < count && int(atomic.LoadInt32(&p.workerCount)) < p.maxWorkers; i++ {
		p.wg.Add(1)
		atomic.AddInt32(&p.workerCount, 1)
		go p.worker()
	}
}

// worker 工作协程
func (p *Pool) worker() {
	defer func() {
		p.wg.Done()
		atomic.AddInt32(&p.workerCount, -1)

		if r := recover(); r != nil {
			log.Printf("worker recovered from panic: %v", r)
		}
	}()

	for {
		select {
		case task := <-p.taskChan:
			if task != nil {
				p.processTask(task)
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// processTask 处理任务
func (p *Pool) processTask(task *Task) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("task processing panic: %v", r)
			atomic.AddInt64(&p.stats.FailedTasks, 1)
		}
	}()

	atomic.AddInt64(&p.stats.ProcessedTasks, 1)

	switch task.Type {
	case TaskTypePage:
		p.processPageTask(task)
	case TaskTypeItem:
		p.processItemTask(task)
	}
}

// processPageTask 处理页面任务
func (p *Pool) processPageTask(task *Task) {
	log.Printf("Processing page %d: %s", task.PageNum, task.URL)

	doc := GetHtmlNode(task.URL)
	if doc == nil {
		log.Printf("Failed to get HTML for page %d", task.PageNum)
		return
	}

	// 解析列表项
	ch := make(chan *html.Node)
	go func() {
		ParseNode(doc, &model.Node{
			Type:      html.ElementNode,
			Data:      "li",
			ClassName: "clearfloat",
		}, ch)
	}()

	// 为每个列表项创建新任务
	for n := range ch {
		itemTask := &Task{
			Type: TaskTypeItem,
			Data: n,
			URL:  task.URL, // 传递基础URL用于构建完整链接
		}
		p.Submit(itemTask)
	}
}

// processItemTask 处理列表项任务
func (p *Pool) processItemTask(task *Task) {
	node, ok := task.Data.(*html.Node)
	if !ok {
		log.Printf("Invalid task data type")
		return
	}

	res := &model.NotiInfo{}
	var href string

	// 解析列表项内容
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "p" {
			for nc := c.FirstChild; nc != nil; nc = nc.NextSibling {
				if nc.Type != html.ElementNode {
					continue
				}
				if nc.Data == "a" {
					if util.GetHtmlNodeValByKey(nc, "class") == "lm_a" {
						res.Author = strings.TrimSpace(nc.FirstChild.Data)
					} else {
						res.Title = util.GetHtmlNodeValByKey(nc, "title")
						href = util.GetHtmlNodeValByKey(nc, "href")
					}
				} else {
					res.Time = strings.TrimSpace(nc.FirstChild.Data)
				}
			}
		}
	}

	// 时间过滤
	if res.Time < "2020-01-01" || res.Time > "2021-09-01" {
		return
	}

	// 获取详细信息
	detailURL := os.Getenv("BASE_URL") + href
	doc := GetHtmlNode(detailURL)
	if doc == nil {
		log.Printf("Failed to get detail HTML: %s", detailURL)
		return
	}

	ch := make(chan *html.Node)
	go func() {
		ParseNode(doc, &model.Node{
			Type:      html.ElementNode,
			Data:      "div",
			ClassName: "conthsj",
		}, ch)
	}()

	for nn := range ch {
		for c := nn.FirstChild; c != nil; c = c.NextSibling {
			if c.Data != "script" {
				continue
			}
			if clicktype, owner, clickid, ok := util.ParseShowDynClicks(strings.TrimSpace(c.FirstChild.Data)); ok {
				cnt, err := GetFZUClickCount(clicktype, owner, clickid)
				if err != nil {
					log.Printf("获取点击数失败: %v", err)
					break
				}
				res.Count = cnt
			}
			break
		}
	}

	// 保存到数据库
	if err := db.InitDb().AddNi(res); err != nil {
		log.Printf("保存数据失败: %v", err)
	}
}

// Submit 提交任务
func (p *Pool) Submit(task *Task) bool {
	if atomic.LoadInt32(&p.isStopped) == 1 {
		return false
	}

	select {
	case p.taskChan <- task:
		atomic.AddInt64(&p.stats.TotalTasks, 1)
		return true
	default:
		log.Printf("Task queue is full, consider increasing queue size")
		return false
	}
}

// monitor 监控协程
func (p *Pool) monitor() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentWorkers := atomic.LoadInt32(&p.workerCount)
			queueLen := len(p.taskChan)

			log.Printf("Pool stats - Workers: %d, Queue: %d, Processed: %d/%d",
				currentWorkers, queueLen,
				atomic.LoadInt64(&p.stats.ProcessedTasks),
				atomic.LoadInt64(&p.stats.TotalTasks))

			// 动态调整worker数量
			if queueLen > int(currentWorkers)*3 && int(currentWorkers) < p.maxWorkers {
				p.startWorkers(2)
				log.Printf("Added 2 workers due to queue backlog")
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// Stop 停止协程池
func (p *Pool) Stop() {
	if !atomic.CompareAndSwapInt32(&p.isStopped, 0, 1) {
		return
	}

	p.cancel()
	close(p.taskChan)
	p.wg.Wait()

	log.Printf("Pool stopped. Final stats - Processed: %d, Failed: %d",
		atomic.LoadInt64(&p.stats.ProcessedTasks),
		atomic.LoadInt64(&p.stats.FailedTasks))
}

// GetStats 获取统计信息
func (p *Pool) GetStats() (int32, int, int64, int64) {
	return atomic.LoadInt32(&p.workerCount),
		len(p.taskChan),
		atomic.LoadInt64(&p.stats.ProcessedTasks),
		atomic.LoadInt64(&p.stats.FailedTasks)
}
