package main

import (
	"spider/engine"
	"spider/parser/huaban"
)

func main() {
	e := engine.New()
	e.WorkerNum = 1000
	req := make([]engine.Request, 0)
	req = append(req, engine.Request{
		Url:  "http://huaban.com/partner/uc/aimeinv/pins/",
		Func: huaban.ParserImage,
	})
	//req = append(req, engine.Request{
	//	Url:  "http://thznn.com/thread-1967472-1-10.html",
	//	Func: taohuazu.ParserImage,
	//})
	e.Run(req...)
}
