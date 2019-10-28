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
		URL:  "http://thznn.com/thread-1967472-1-10.html",
		FUNC: new(parser.TaoHua).Handle,
	})
	e.AddRequest(engine.Request{
		URL:  "http://huaban.com/partner/uc/aimeinv/pins/",
		FUNC: new(parser.HuaBan).Handle,
	})
	e.Start()
}
