package HttpHandler

import (
	"MyHoneyProxy/HoneyProxy"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpMiddle struct {
	RequestBody []byte
	ResponseBody []byte
}

func copyRequestBody(res *http.Request) ([]byte, error) {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return buf, nil
}

func HandleRequest(req *http.Request, ctx *HoneyProxy.ProxyCtx) (*http.Request, *http.Response) {
	var err error
	httpMiddle := HttpMiddle{}
	ctx.UserData = &httpMiddle
	//获取请求内容
	httpMiddle.RequestBody, err = copyRequestBody(req)
	if err != nil {
		return req, nil
	}
	log.Println("执行请求:",req.URL.String())
	//演示如何替换resp
	if req.URL.Host == "baidu.com"{
		return req,&http.Response{}
	}
	return req, nil
}