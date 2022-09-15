package HoneyProxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

type protocol uint8
const(
	socks4 protocol = 0x1
	socks5 protocol = 0x2
	head protocol = 0x3
	delete protocol = 0x4
	connect protocol = 0x5
	options protocol = 0x6
	trace protocol = 0x7
)

var(
	uClient = http.Client{
		Transport:&http.Transport{
			TLSClientConfig:&tls.Config{InsecureSkipVerify:true}},
	}
)

type HoneyProxy struct {
	//参数req为请求,返回req为新请求,返回resp不为空表示填充响应
	reqHandler	func(req *http.Request, ctx *ProxyCtx) (*http.Request, *http.Response)
	//参数resp为响应,返回resp为新响应
	respHandler func(resp *http.Response, ctx *ProxyCtx)*http.Response
	Tr *http.Transport
}

func NewHoneyProxy()*HoneyProxy{
	return &HoneyProxy{
		Tr:&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, Proxy: http.ProxyFromEnvironment},
	}
}

func (this *HoneyProxy)filterRequest(req *http.Request, ctx *ProxyCtx)(*http.Request, *http.Response) {
	if this.reqHandler != nil{
		return this.reqHandler(req,ctx)
	}
	return req,nil
}

func (this *HoneyProxy)filterResponse(respOrig *http.Response, ctx *ProxyCtx)*http.Response  {
	if this.respHandler != nil{
		return this.respHandler(respOrig,ctx)
	}
	return respOrig
}

func (this *HoneyProxy)SetReqHandler(reqHandler func(req *http.Request, ctx *ProxyCtx) (*http.Request, *http.Response))  {
	this.reqHandler = reqHandler
}

func (this *HoneyProxy)SetRespHandler(respHandler func(resp *http.Response, ctx *ProxyCtx)*http.Response)  {
	this.respHandler = respHandler
}

func (this *HoneyProxy)checkProxy()  {



}

func (this *HoneyProxy)handleSocks4Request()error  {





	return nil
}


func (this *HoneyProxy)handleHttpRequest(conn net.Conn,tmpBuffer []byte)error {
	req,err := http.ReadRequest(bufio.NewReader(bytes.NewReader(tmpBuffer)))
	if err != nil{
		return err
	}
	req.RequestURI = ""
	resp,err := uClient.Do(req)
	if err != nil{
		return err
	}
	defer resp.Body.Close()
	respBytes,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return err
	}
	var outBuffer bytes.Buffer
	outBuffer.WriteString("HTTP/1.0")
	outBuffer.WriteByte(0x20)
	outBuffer.WriteString(resp.Status)
	outBuffer.WriteString("\r\n")
	for eKey,eValue := range resp.Header{
		outBuffer.WriteString(eKey)
		outBuffer.WriteString(": ")
		if len(eValue) > 0 {
			outBuffer.WriteString(eValue[0])
		}
		outBuffer.WriteString("\r\n")
	}
	conn.Write(outBuffer.Bytes())
	conn.Write([]byte("\r\n"))
	conn.Write(respBytes)
	return nil
}

func (this *HoneyProxy)handleSocks5Request()error  {

	return nil
}

//接收完整的数据,以\r\n\r\n作为结尾
func (this *HoneyProxy)readCompleteBuffer(conn net.Conn)(retBuf []byte,retErr error)  {
	for true{
		tmpBuffer := make([]byte,0x2000)
		nLen, err := conn.Read(tmpBuffer)
		if err != nil{
			return retBuf,err
		}
		retBuf = append(retBuf, tmpBuffer[0:nLen]...)
		if bytes.HasSuffix(retBuf,[]byte{0xD,0xA,0xD,0xA}) == true{
			return retBuf,nil
		}
	}
	return retBuf,nil
}

func (this *HoneyProxy)handleConn(conn net.Conn)error  {
	defer conn.Close()
	//读取完整的内容
	tmpBuffer,err := this.readCompleteBuffer(conn)
	if err != nil{
		return err
	}
	switch tmpBuffer[0] {
	case 0x4:
		return this.handleSocks4Request()
	case 0x5:
		return this.handleSocks5Request()
	case 'O':	//options
		return this.handleHttpRequest(conn,tmpBuffer)
	case 'P':	//post
		return this.handleHttpRequest(conn,tmpBuffer)
	case 'T':	//trace
		return this.handleHttpRequest(conn,tmpBuffer)
	case 'D':	//delete
		return this.handleHttpRequest(conn,tmpBuffer)
	case 'G':	//get
		return this.handleHttpRequest(conn,tmpBuffer)
	case 'C':	//connect
		return this.handleHttpsRequest(conn,tmpBuffer)
	}
	return nil
}

func (this *HoneyProxy)StartProxy(port int,maxUser int)error  {
	ls,err := net.Listen("tcp",":" + strconv.Itoa(port))
	if err != nil{
		return err
	}
	defer ls.Close()
	for true{
		conn,err := ls.Accept()
		if err != nil{
			continue
		}
		go this.handleConn(conn)
	}
	return nil
}
