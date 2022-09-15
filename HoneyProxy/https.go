package HoneyProxy

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var(
	httpsRegexp = regexp.MustCompile(`^https:\/\/`)
)

//处理真正的https请求

func (this *HoneyProxy)executeHttpsRequest(clientTls net.Conn,ctx *ProxyCtx)error  {
	var err error
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

func (this *HoneyProxy)handleHttpsRequest(proxyClient net.Conn,proxyReq *http.Request,ctx *ProxyCtx)error  {
	proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	targetSiteCon,err := net.Dial("tcp",proxyReq.Host)
	if err != nil{
		return err
	}
	defer targetSiteCon.Close()
	ctx.ParseProxyAuth(proxyReq)
	tlsConfig, err := TLSConfigFromCA()(proxyReq.Host,ctx)
	if err != nil{
		return err
	}
	rawClientTls := tls.Server(proxyClient,tlsConfig)
	defer rawClientTls.Close()
	err = rawClientTls.Handshake()
	if err != nil {
		return err
	}
	cReq,err := http.ReadRequest(bufio.NewReader(rawClientTls))
	if err != nil{
		return err
	}
	//开始整理真正的请求头
	cReq.RemoteAddr = proxyClient.RemoteAddr().String()
	if !httpsRegexp.MatchString(cReq.URL.String()) {
		cReq.URL, err = url.Parse("https://" + cReq.Host + cReq.URL.String())
	}
	ctx.Req = cReq
	return this.executeHttpsRequest(proxyClient,ctx)
}