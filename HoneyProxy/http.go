package HoneyProxy

import (
	"io"
	"net/http"
	"net/url"
)

func (this *HoneyProxy)handleHttpRequest(bufConn *bufferedConn,ctx *ProxyCtx)error {
	req,err := http.ReadRequest(bufConn.r)
	if err != nil{
		return err
	}
	req.RequestURI = ""
	if ctx.Protocol == protocol_socks4 || ctx.Protocol == protocol_socks5{
		req.URL, err = url.Parse("http://" + req.Host + req.URL.Path)
		if err != nil{
			return err
		}
	}
	ctx.parseBasicAuth(req)
	var resp *http.Response
	req,resp = this.filterRequest(req,ctx)
	if resp == nil{
		resp,err = ctx.RoundTrip(req)
		if err != nil{
			return err
		}
	}
	defer resp.Body.Close()
	_,err = io.WriteString(bufConn,"HTTP/1.0 " + resp.Status + "\r\n")
	if err != nil{
		return err
	}
	err = resp.Header.Write(bufConn)
	if err != nil{
		return err
	}
	bufConn.Write([]byte("\r\n"))
	_,err = io.Copy(bufConn, resp.Body)
	if err != nil{
		return err
	}
	return nil
}