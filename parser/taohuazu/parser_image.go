package taohuazu

import (
	"spider/engine"
	"github.com/PuerkitoBio/goquery"
	"spider/memory"
)

var (
	host = "http://thznn.com/"
)

func ParserImage(url string) engine.Response {
	response := engine.Response{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		response.Requests = append(response.Requests, engine.Request{
			Url:  url,
			Func: ParserImage,
		})
		return response
	}

	doc.Find("div.t_fsz > table > tbody > tr > td.t_f > img").Each(func(i int, s *goquery.Selection) {
		if imgUrl, e := s.Attr("file"); e {
			// 处理前缀为//的url
			if string([]rune(imgUrl)[:2]) == "//" {
				imgUrl = "http:" + imgUrl
			}
			// 无Host的url
			if string([]rune(imgUrl)[:4]) != "http" {
				imgUrl = host + imgUrl
			}
			response.Results = append(response.Results, engine.Result{
				Addr: imgUrl,
				Func: memory.Download,
			})
		}
	})

	nUrl, e := doc.Find(".pcb > a").Eq(1).Attr("href")
	if !e {
		return response
	}
	// 处理前缀为//的url
	if string([]rune(nUrl)[:2]) == "//" {
		nUrl = "http:" + nUrl
	}
	// 无Host的url
	if string([]rune(nUrl)[:4]) != "http" {
		nUrl = host + nUrl
	}
	response.Requests = append(response.Requests, engine.Request{
		Url:  nUrl,
		Func: ParserImage,
	})
	return response
}
