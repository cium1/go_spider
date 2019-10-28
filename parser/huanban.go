package parser

import (
	"go_spider/engine"
	"regexp"
	"strings"
	"time"
)

type HuaBan struct{}

var (
	//基础URL
	hHost = "http://hbimg.b0.upaiyun.com/"
	hBase = "http://huaban.com/partner/uc/aimeinv/pins/"
)

var (
	//正则编译
	hKeyMust = regexp.MustCompile(`"key":"(.*?)"`)
	hIdMust  = regexp.MustCompile(`"pin_id":(\d+),`)
)

var (
	//定时器限速
	timer = time.Tick(time.Millisecond * 10)
)

func (h *HuaBan) Handle(url string) (scheduler engine.Scheduler) {

	//<-timer

	scheduler = engine.Scheduler{}

	body, err := engine.GetBody(url)
	if err != nil {
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			URL:  url,
			FUNC: h.Handle,
		})
		return
	}

	images := hKeyMust.FindAllStringSubmatch(string(body), -1)
	for _, val := range images {
		addr := val[1]
		//过滤图标文件
		if !strings.Contains(addr, "-") {
			continue
		}
		//修正URL
		if string([]rune(addr)[:4]) != "http" {
			addr = hHost + addr
		}
		scheduler.Processors = append(scheduler.Processors, engine.Processor{
			Content: addr,
			FUNC:    engine.Download,
		})
	}

	ids := hIdMust.FindAllStringSubmatch(string(body), -1)
	if len(ids) > 0 {
		id := ids[len(ids)-1][1]
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			URL:  hBase + "?max=" + id + "&limit=8&wfl=1",
			FUNC: h.Handle,
		})
	}

	return
}
