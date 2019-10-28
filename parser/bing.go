package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go_spider/engine"
	"strings"
)

type Bing struct {
}

func (b *Bing) Home(url string) (scheduler engine.Scheduler) {

	scheduler = engine.Scheduler{}

	body, err := engine.GetBody(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		if url, e := selection.Attr("href"); e {
			if !strings.Contains(url, "http") {
				url = "https:" + url
			}
			scheduler.Requests = append(scheduler.Requests, engine.Request{
				URL:  url,
				FUNC: b.Home,
			})
		}
	})

	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		if img, e := selection.Attr("src"); e {
			fmt.Println(img)
		}
	})

	return
}
