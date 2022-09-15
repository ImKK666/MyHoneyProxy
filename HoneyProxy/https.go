package HoneyProxy

import (
	"bufio"
	"bytes"
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

func (this *HoneyProxy)handleHttpsRequest(proxyClient net.Conn,tmpBuffer []byte)error  {
	proxyReq,err := http.ReadRequest(bufio.NewReader(bytes.NewReader(tmpBuffer)))
	if err != nil{
		return err
	}
	proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	targetSiteCon,err := net.Dial("tcp",proxyReq.Host)
	if err != nil{
		return err
	}
	defer targetSiteCon.Close()
	ctx := &ProxyCtx{Proxy: this}
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
	var resp *http.Response
	cReq,resp = this.filterRequest(cReq,ctx)
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
	_, err = io.WriteString(rawClientTls, "HTTP/1.1" + " " + statusCode + text + "\r\n")
	if err != nil{
		return err
	}
	resp.Header.Del("Content-Length")
	resp.Header.Set("Transfer-Encoding", "chunked")
	// Force connection close otherwise chrome will keep CONNECT tunnel open forever
	resp.Header.Set("Connection", "close")
	err = resp.Header.Write(rawClientTls)
	if err != nil {
		return err
	}
	_, err = io.WriteString(rawClientTls, "\r\n")
	if err != nil {
		return err
	}
	chunked := newChunkedWriter(rawClientTls)
	_, err = io.Copy(chunked, resp.Body)
	if err != nil {
		return err
	}
	err = chunked.Close()
	if err != nil {
		return err
	}
	_, err = io.WriteString(rawClientTls, "\r\n")
	if err != nil {
		return err
	}
	return nil
}