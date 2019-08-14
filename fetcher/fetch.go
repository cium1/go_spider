package fetcher

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
)

var (
//定时器
//timer = time.Tick(time.Millisecond * 10)
)

// http请求
func Fetch(url string) ([]byte, error) {

	//<-timer

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	res, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(res)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error : %v %v", url, resp.StatusCode)
	}

	respBody := bufio.NewReader(resp.Body)
	utf8Reader := transform.NewReader(respBody, DetermineEncoding(respBody).NewDecoder())
	body, err := ioutil.ReadAll(utf8Reader)

	return body, err
}

// 识别字符编码
func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	data, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(data, "")
	return e
}
