package parser

import (
	"github.com/PuerkitoBio/goquery"
	"go_spider/engine"
	"net/http"
)

type TaoHua struct {
}

var (
	tHost = "http://xxx.com/"
)

func (t *TaoHua) Handle(url string) (scheduler engine.Scheduler) {

	scheduler = engine.Scheduler{}

	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			URL:  url,
			FUNC: t.Handle,
		})
		return
	}

	doc.Find("div.t_fsz > table > tbody > tr > td.t_f > img").Each(func(i int, s *goquery.Selection) {

		if imgUrl, e := s.Attr("file"); e {

			if string([]rune(imgUrl)[:2]) == "//" {
				imgUrl = "http:" + imgUrl
			}

			if string([]rune(imgUrl)[:4]) != "http" {
				imgUrl = tHost + imgUrl
			}

			scheduler.Processors = append(scheduler.Processors, engine.Processor{
				Content: imgUrl,
				FUNC:    engine.Download,
			})
		}

	})

	nUrl, e := doc.Find(".pcb > a").Eq(1).Attr("href")
	if !e {
		return
	}

	if string([]rune(nUrl)[:2]) == "//" {
		nUrl = "http:" + nUrl
	}

	if string([]rune(nUrl)[:4]) != "http" {
		nUrl = tHost + nUrl
	}

	scheduler.Requests = append(scheduler.Requests, engine.Request{
		URL:  nUrl,
		FUNC: t.Handle,
	})

	return
}
