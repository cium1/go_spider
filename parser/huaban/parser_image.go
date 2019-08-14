package huaban

import (
	"regexp"
	"spider/engine"
	"spider/fetcher"
	"spider/memory"
)

const (
	host = "http://hbimg.b0.upaiyun.com/"
	base = "http://huaban.com/partner/uc/aimeinv/pins/"
)

var (
	keyMust = regexp.MustCompile(`"key":"(.*?)"`)
	idMust  = regexp.MustCompile(`"pin_id":(\d+),`)
)

func ParserImage(url string) engine.Response {
	response := engine.Response{}
	//获取响应
	body, err := fetcher.Fetch(url)
	if err != nil {
		response.Requests = append(response.Requests, engine.Request{
			Url:  url,
			Func: ParserImage,
		})
		return response
	}

	//解析图片
	images := keyMust.FindAllStringSubmatch(string(body), -1)
	for _, val := range images {
		addr := val[1]
		// 处理前缀为//的url
		if string([]rune(addr)[:2]) == "//" {
			addr = "http:" + addr
		}
		// 无Host的url
		if string([]rune(addr)[:4]) != "http" {
			addr = host + addr
		}
		response.Results = append(response.Results, engine.Result{
			Addr: addr,
			Func: memory.Download,
		})
	}

	//解析下一页
	ids := idMust.FindAllStringSubmatch(string(body), -1)
	if len(ids) > 0 {
		id := ids[len(ids)-1][1]
		response.Requests = append(response.Requests, engine.Request{
			Url:  base + "?max=" + id + "&limit=8&wfl=1",
			Func: ParserImage,
		})
	}

	return response
}
