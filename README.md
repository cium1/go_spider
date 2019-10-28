### Go 爬虫引擎 + 实例(parser)
#### 仅用于娱乐与学习

```go

    func main() {
        runtime.GOMAXPROCS(runtime.NumCPU())
        e := engine.New()
        e.AddRequest(engine.Request{
            URL:  "http://xxx.com/thread-1967472-1-10.html",
            FUNC: new(parser.TaoHua).Handle,
        }, engine.Request{
            URL:  "http://huaban.com/partner/uc/aimeinv/pins/",
            FUNC: new(parser.HuaBan).Handle,
        })
        e.Start()
    }

```