package HoneyProxy

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//处理真正的https请求

func (this *HoneyProxy)executeHttpsRequest(clientTls net.Conn,ctx *ProxyCtx)error  {

	cReq,err := http.ReadRequest(bufio.NewReader(clientTls))
	if err != nil{
		return err
	}

	//整理真正的请求头
	if strings.HasPrefix(cReq.URL.String(),"https://") == false {
		cReq.URL, err = url.Parse("https://" + cReq.Host + cReq.URL.String())
	}

	ctx.Req = cReq
	cReq,resp := this.filterRequest(ctx.Req,ctx)
	if resp == nil{
		removeProxyHeaders(cReq)
		resp, err = ctx.RoundTrip(cReq)
		if err != nil{
			return err
		}
	}
	defer resp.Body.Close()
	resp = this.filterResponse(resp,ctx)
	if resp == nil{
		return err
	}
	text := resp.Status
	statusCode := strconv.Itoa(resp.StatusCode) + " "
	if strings.HasPrefix(text, statusCode) {
		text = text[len(statusCode):]
	}
	resp = this.filterResponse(resp, ctx)
	_, err = io.WriteString(clientTls, "HTTP/1.1" + " " + statusCode + text + "\r\n")
	if err != nil{
		return err
	}
	resp.Header.Del("Content-Length")
	resp.Header.Set("Transfer-Encoding", "chunked")
	// Force connection close otherwise chrome will keep CONNECT tunnel open forever
	resp.Header.Set("Connection", "close")
	err = resp.Header.Write(clientTls)
	if err != nil {
		return err
	}
	_, err = io.WriteString(clientTls, "\r\n")
	if err != nil {
		return err
	}
	chunked := newChunkedWriter(clientTls)
	_, err = io.Copy(chunked, resp.Body)
	if err != nil {
		return err
	}
	err = chunked.Close()
	if err != nil {
		return err
	}
	_, err = io.WriteString(clientTls, "\r\n")
	if err != nil {
		return err
	}
	return nil
}

func (this *HoneyProxy)https_prepareRequest(proxyClient *bufferedConn,ctx *ProxyCtx)(*tls.Config,error)  {
	req,err := http.ReadRequest(proxyClient.r)
	if err != nil{
		return nil,err
	}
	ctx.parseBasicAuth(req)
	_,err = proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	if err != nil{
		return nil,err
	}
	tlsConfig, err := TLSConfigFromCA(req.Host,ctx)
	if err != nil{
		return nil,err
	}
	return tlsConfig,nil
}

func (this *HoneyProxy)handleHttpsRequest(proxyClient *bufferedConn,tlsConfig *tls.Config,ctx *ProxyCtx)error  {
	rawClientTls := tls.Server(proxyClient,tlsConfig)
	defer rawClientTls.Close()
	err := rawClientTls.Handshake()
	if err != nil {
		return err
	}
	return this.executeHttpsRequest(rawClientTls,ctx)
}