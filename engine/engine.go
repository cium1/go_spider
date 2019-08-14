package engine

type Engine struct {
	WorkerNum int
}

func New() *Engine {
	return new(Engine)
}

func (e *Engine) Run(seeds ...Request) {

	req := make(chan Request, 100)
	resp := make(chan Response, 100)
	res := make(chan Result, 100)

	//创建worker
	for i := 0; i < e.WorkerNum; i++ {
		go e.worker(req, resp, res)
	}

	//发送初始请求
	for _, r := range seeds {
		req <- r
	}

	select {}
}

func (e *Engine) worker(req chan Request, resp chan Response, res chan Result) {

	for {
		select {

		case request := <-req:
			resp <- request.Func(request.Url)

		case response := <-resp:
			for _, request := range response.Requests {
				req <- request
			}
			for _, result := range response.Results {
				res <- result
			}

		case result := <-res:
			resp <- result.Func(result.Addr)

		}
	}

}
