package HttpHandler

import (
	"MyHoneyProxy/HoneyProxy"
	"MyHoneyProxy/Manager/UploadManager"
	"MyHoneyProxy/Model"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var(
	//用来执行过滤的ContentTYpe
	gContentTypeMap = map[string]struct{}{
		"application/octet-stream":{},
		"application/x-rar-compressed":{},
		"application/zip":{},
		"application/x-javascript":{},
		"application/ogg":{},
		"application/x-shockwave-flash":{},
		"image/png": {},
		"image/jpeg":{},
		"image/gif":{},
		"image/x-icon":{},
	}
)

func copyResponse(resp *http.Response)[]byte {
	data, _ := io.ReadAll(resp.Body)
	if len(data) == 0{
		return nil
	}
	resp.Body = io.NopCloser(bytes.NewReader(data))
	return data
}

func pushHoneyData(honeyReq *Model.HoneyReq)error  {
	//限制响应头长度为1MB
	if len(honeyReq.ResponseBody) > 1048576 {
		honeyReq.ResponseBody = honeyReq.ResponseBody[0:1048576]
	}
	honeyBytes,err := json.Marshal(honeyReq)
	if err != nil{
		return err
	}
	UploadManager.Instance.PushHoneyData(honeyBytes)
	return nil
}

func HandleResponse(resp *http.Response, ctx *HoneyProxy.ProxyCtx) *http.Response {
	if ctx.UserData == nil {
		return resp
	}
	httpMiddle := ctx.UserData.(*Model.HttpMiddle)
	if httpMiddle.SkipBody == false{
		_,bFilter := gContentTypeMap[resp.Header.Get("Content-Type")]
		if bFilter == false{
			httpMiddle.ResponseBody = copyResponse(resp)
		}
	}
	ctx.Req.URL.Host = httpMiddle.Host
	//tmpReq := Model.HoneyReq{
	//	Method:ctx.Req.Method,
	//	RemoteAddress:httpMiddle.IpAddr,
	//	Scheme: ctx.Req.URL.Scheme,
	//	Host:httpMiddle.Host,
	//	Port: httpMiddle.Port,
	//	Path:ctx.Req.URL.Path,
	//	Url:ctx.Req.URL.String(),
	//	RequestParam:ctx.Req.URL.Query(),
	//	RequestHeader:ctx.Req.Header,
	//	RequestBody:httpMiddle.RequestBody,
	//	ResponseHeader:resp.Header,
	//	ResponseBody:httpMiddle.ResponseBody,
	//	Ip:DialManager.Instance.GetIpAddress(),
	//	Origin:ConfigManager.Instance.ProxyServer.Origin,
	//	OriginDetail:ConfigManager.Instance.ProxyServer.OriginDetail,
	//	UserName: ctx.ProxyAuth.UserName,
	//	Password:ctx.ProxyAuth.PassWord,
	//	StatusCode:resp.StatusCode,
	//	AttackTime:time.Now().Unix(),
	//}
	//err = pushHoneyData(&tmpReq)

	return resp
}