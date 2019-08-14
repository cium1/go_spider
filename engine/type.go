package engine

type Request struct {
	Url  string
	Func func(url string) Response
}

type Response struct {
	Requests []Request
	Results  []Result
}

type Result struct {
	Addr string
	Func func(addr string) Response
}
