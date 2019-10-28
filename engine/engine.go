package engine

import (
	"fmt"
	"sync"
	"time"
)

type Request struct {
	URL  string
	FUNC func(url string) Scheduler
}

type Scheduler struct {
	Requests   []Request
	Processors []Processor
}

type Processor struct {
	Content string
	FUNC    func(content string) Scheduler
}

type Engine struct {
	WorkerNum int
	Requests  []Request
}

var (
	wg sync.WaitGroup
)

// 初始化
func New() *Engine {
	return new(Engine)
}

// 初始请求
func (e *Engine) AddRequest(request ...Request) {
	e.Requests = append(e.Requests, request...)
}

// 开始
func (e *Engine) Start() {

	start := time.Now()

	if e.WorkerNum == 0 {
		e.WorkerNum = 1000
	}

	request, response, result := make(chan Request, e.WorkerNum), make(chan Scheduler, e.WorkerNum), make(chan Processor, e.WorkerNum)

	//发送初始请求
	for _, r := range e.Requests {
		request <- r
	}

	//创建worker
	for i := 0; i < e.WorkerNum; i++ {
		wg.Add(1)
		go e.worker(request, response, result)
	}

	wg.Wait()

	d, _ := time.ParseDuration("60s")
	fmt.Println(time.Since(start.Add(d)))
}

// 创建进程
func (e *Engine) worker(requestChan chan Request, responseChan chan Scheduler, resultChan chan Processor) {
	defer wg.Done()
	for {

		select {

		case request := <-requestChan:
			if request.FUNC != nil {
				responseChan <- request.FUNC(request.URL)
			}

		case response := <-responseChan:
			if response.Requests != nil {
				for _, request := range response.Requests {
					requestChan <- request
				}
			}
			if response.Processors != nil {
				for _, result := range response.Processors {
					resultChan <- result
				}
			}

		case result := <-resultChan:
			if result.FUNC != nil {
				responseChan <- result.FUNC(result.Content)
			}

		case <-time.After(time.Duration(time.Second * 60)): //超时
			return

		}
	}

}
