package HoneyProxy

import "net/http"

type ProxyAuth struct {
	UserName string
	PassWord string
}

type RoundTripper interface {
	RoundTrip(req *http.Request, ctx *ProxyCtx) (*http.Response, error)
}

type ProxyCtx struct {
	//协议类型
	Protocol protocol
	Req *http.Request
	ProxyAuth ProxyAuth
	RoundTripper RoundTripper
	UserData interface{}
	Proxy     *HoneyProxy
}

func (ctx *ProxyCtx) RoundTrip(req *http.Request) (*http.Response, error) {
	if ctx.RoundTripper != nil {
		return ctx.RoundTripper.RoundTrip(req, ctx)
	}
	return ctx.Proxy.Tr.RoundTrip(req)
}
