package engine

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Scheduler struct {
	Requests   []Request
	Processors []Processor
}

type Engine struct {
	WorkerNum int
	WorkerI   int
	TimeOut   time.Duration
	Requests  []Request
}

type Request struct {
	Url     string
	Operate interface{}
	Func    func(url string, operate interface{}) Scheduler
}

type Processor struct {
	Content string
	Operate interface{}
	Func    func(content string, operate interface{}) Scheduler
}

var (
	wg   sync.WaitGroup
	lock sync.Mutex
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
		e.WorkerNum = runtime.NumCPU() * 2
	}

	if e.TimeOut == 0 {
		e.TimeOut = time.Minute * 1
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

	fmt.Println("运行时长", time.Since(start.Add(e.TimeOut)))
}

// 创建进程
func (e *Engine) worker(requestChan chan Request, responseChan chan Scheduler, resultChan chan Processor) {
	defer wg.Done()
	for {

		select {

		case request := <-requestChan:
			fmt.Print(".")
			if request.Func != nil {
				responseChan <- request.Func(request.Url, request.Operate)
			}

		case response := <-responseChan:
			fmt.Print(".")
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
			fmt.Print(".")
			if result.Func != nil {
				responseChan <- result.Func(result.Content, result.Operate)
			}

		case <-time.After(e.TimeOut): //超时
			lock.Lock()
			e.WorkerI++
			fmt.Println("timeout", e.WorkerI)
			lock.Unlock()
			return

		}
	}

}
