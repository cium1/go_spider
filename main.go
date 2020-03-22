package main

import (
	"go_spider/engine"
	"go_spider/parser"
	"time"
)

func main() {
	e := engine.New()
	e.WorkerNum = 200
	e.TimeOut = time.Minute * 2
	e.AddRequest(engine.Request{
		Url:  "http://xxx.cc/thread-1930952-1-1.html",
		Func: (&parser.TaoHua{Dir: "../taohuazhu"}).Request,
	})
	e.Start()
}
