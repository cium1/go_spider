package parser

import (
	"bufio"
	"fmt"
	"go_spider/engine"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// TaoHua is struct
type TaoHua struct {
	Dir string
}

// Download is images
func (t *TaoHua) Download(content string, operate interface{}) (scheduler engine.Scheduler) {
	scheduler = engine.Scheduler{}

	if t.Dir == "" {
		log.Println("No directory option exists dir")
		return
	}

	path := t.Dir + string(os.PathSeparator) + operate.(string) + string(os.PathSeparator)
	if err := Mkdir(path); err != nil {
		log.Println("mkdir", t.Dir, err)
		return
	}

	fileName := filepath.Base(content)

	response, err := HTTPGet(content)
	if err != nil {
		log.Println(err)
		if err == http.ErrHandlerTimeout {
			scheduler.Processors = append(scheduler.Processors, engine.Processor{
				Content: content,
				Operate: operate,
				Func:    t.Download,
			})
		}
		return
	}
	defer response.Body.Close()

	file, err := os.Create(filepath.Dir(path) + string(os.PathSeparator) + fileName)
	if err != nil {
		log.Println(err)
		scheduler.Processors = append(scheduler.Processors, engine.Processor{
			Content: content,
			Operate: operate,
			Func:    t.Download,
		})
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	reader := bufio.NewReaderSize(response.Body, 32*1024)

	size, err := io.Copy(writer, reader)
	if err != nil {
		log.Println(err)
		scheduler.Processors = append(scheduler.Processors, engine.Processor{
			Content: content,
			Operate: operate,
			Func:    t.Download,
		})
		return
	}
	writer.Flush()

	fmt.Println(content, engine.FileByte(engine.FileSize(size)))

	return
}

// WriteFile is content to file
func (t *TaoHua) WriteFile(content string, operate interface{}) (scheduler engine.Scheduler) {
	scheduler = engine.Scheduler{}

	if t.Dir == "" {
		log.Println("No directory option exists dir")
		return
	}

	path := t.Dir + string(os.PathSeparator) + operate.(string) + string(os.PathSeparator)
	if err := Mkdir(path); err != nil {
		log.Println("mkdir", t.Dir, err)
		return
	}

	file, err := os.Create(filepath.Dir(path) + string(os.PathSeparator) + "link.html")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	_, _ = file.WriteString(fmt.Sprintf("<script>window.location.href = '%s'</script>", content))

	return
}

// ZZ is 种子解析
func (t *TaoHua) ZZ(url string, operate interface{}) (scheduler engine.Scheduler) {

	scheduler = engine.Scheduler{}

	response, err := HTTPGet(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println(err)
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			Url:     url,
			Operate: operate,
			Func:    t.ZZ,
		})
		return
	}

	file, _ := doc.Find("div#wp > div.f_c > div > div > a").Eq(0).Attr("href")

	scheduler.Processors = append(scheduler.Processors, engine.Processor{
		Content: file,
		Operate: operate,
		Func:    t.WriteFile,
	})

	return
}

// Request is 发起
func (t *TaoHua) Request(url string, operate interface{}) (scheduler engine.Scheduler) {

	scheduler = engine.Scheduler{}

	response, err := HTTPGet(url)
	if err != nil {
		log.Println(err)
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			Url:  url,
			Func: t.Request,
		})
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println(err)
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			Url:  url,
			Func: t.Request,
		})
		return
	}

	title := doc.Find("span#thread_subject").Text()

	// 图片
	doc.Find("div.t_fsz > table > tbody > tr > td.t_f").Eq(0).Find("img").Each(func(i int, selection *goquery.Selection) {
		if file, bo := selection.Attr("file"); bo && len(file) > 0 {
			scheduler.Processors = append(scheduler.Processors, engine.Processor{
				Content: MakeURL(url, file),
				Operate: title,
				Func:    t.Download,
			})
		}
	})

	// 种子
	if zzURL, bo := doc.Find("td.plc > div.pct > div.pcb > div.t_fsz > div.pattl a").Eq(0).Attr("href"); bo && len(zzURL) > 0 {
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			Url:     MakeURL(url, zzURL),
			Operate: title,
			Func:    t.ZZ,
		})
	}

	// 下一个
	if nextURL, bo := doc.Find("td.plc > div.pct > div.pcb > a").Eq(1).Attr("href"); bo && len(nextURL) > 0 {
		scheduler.Requests = append(scheduler.Requests, engine.Request{
			Url:  MakeURL(url, nextURL),
			Func: t.Request,
		})
	}

	return
}

// HTTPGet is send get
func HTTPGet(url string) (*http.Response, error) {

	client := &http.Client{Timeout: time.Duration(time.Minute * 1)}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}

// MakeURL is make url
func MakeURL(oldURL, URL string) string {

	if strings.HasPrefix(URL, "http") {
		return URL
	}

	return fmt.Sprintf("%s/%s", GetDomain(oldURL), strings.Trim(URL, "/"))
}

// GetDomain is url domain
func GetDomain(url string) string {
	var ht string
	if string([]rune(url)[:4]) == "http" {
		ht = strings.TrimRight(strings.Split(url, "//")[0], ":")
		url = strings.Split(url, "//")[1]
	} else {
		ht = "http"
	}
	return fmt.Sprintf("%s://%s", ht, strings.Split(url, "/")[0])
}

// Mkdir is create dir
func Mkdir(dir string) error {
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}
