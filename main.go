package main

import (
	"go_spider/engine"
	"go_spider/parser"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	e := engine.New()
	e.AddRequest(engine.Request{
		URL:  "http://xxx.com/thread-1967472-1-10.html",
		FUNC: (&parser.TaoHua{}).Handle,
	}, engine.Request{
		URL:  "http://huaban.com/partner/uc/aimeinv/pins/",
		FUNC: (&parser.HuaBan{}).Handle,
	})
	e.Start()
}
