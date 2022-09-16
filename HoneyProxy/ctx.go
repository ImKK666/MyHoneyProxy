package HoneyProxy

import (
	"encoding/base64"
	"net"
	"net/http"
	"strings"
)

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
	//远程连接地址
	RemoteAddr net.Addr
	//真正的请求
	Req *http.Request
	ProxyAuth ProxyAuth
	RoundTripper RoundTripper
	UserData interface{}
	Proxy     *HoneyProxy
}

func (ctx *ProxyCtx)parseBasicAuth(req *http.Request)  {
	proxyHeader := req.Header.Get("Proxy-Authorization")
	if proxyHeader == ""{
		return
	}
	const prefix = "Basic "
	if len(proxyHeader) < len(prefix) {
		return
	}
	decodeBytes,err := base64.StdEncoding.DecodeString(proxyHeader[len(prefix):])
	if err != nil{
		return
	}
	cs := string(decodeBytes)
	s := strings.IndexByte(cs,':')
	if s < 0{
		return
	}
	ctx.ProxyAuth.UserName = cs[:s]
	ctx.ProxyAuth.PassWord = cs[s+1:]
	return
}

func (ctx *ProxyCtx)RoundTrip(req *http.Request) (*http.Response, error) {
	if ctx.RoundTripper != nil {
		return ctx.RoundTripper.RoundTrip(req, ctx)
	}
	return ctx.Proxy.Tr.RoundTrip(req)
}
