package memory

import (
	"path/filepath"
	"github.com/satori/go.uuid"
	"net/http"
	"bufio"
	"io"
	"os"
	"crypto/tls"
	"spider/engine"
	"fmt"
)

var (
	dir = "images/"
)

func Download(addr string) engine.Response {
	response := engine.Response{}
	//初始化存储目录
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	//获取文件后缀名称
	ext := filepath.Ext(filepath.Base(addr))
	// 解决文件无图片格式后缀的问题
	if len(ext) <= 0 {
		ext = ".png"
	}
	//生成唯一文件名称
	fileName, _ := uuid.NewV1()

	//读取文件
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return response
	}
	resp, err := client.Do(res)
	if err != nil {
		//response.Results = append(response.Results, engine.Result{
		//	Addr: addr,
		//	Func: Download,
		//})
		fmt.Println("error:", err)
		return response
	}
	defer resp.Body.Close()

	//创建文件
	reader := bufio.NewReaderSize(resp.Body, 64*1024)
	file, err := os.Create(filepath.Dir(dir) + string(os.PathSeparator) + fileName.String() + ext)
	if err != nil {
		//response.Results = append(response.Results, engine.Result{
		//	Addr: addr,
		//	Func: Download,
		//})
		fmt.Println("error:", err)
		return response
	}

	//写入文件
	_, err = io.Copy(file, reader)
	if err != nil {
		//response.Results = append(response.Results, engine.Result{
		//	Addr: addr,
		//	Func: Download,
		//})
		fmt.Println("error:", err)
		return response
	}
	fmt.Printf("\r%s", addr)
	return response
}
