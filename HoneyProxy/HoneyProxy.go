package HoneyProxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
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

}

func NewHoneyProxy()*HoneyProxy{
	return &HoneyProxy{

	}
}

func (this *HoneyProxy)SetReqHandler()  {
	
}

func (this *HoneyProxy)SetRespHandler()  {


}

func (this *HoneyProxy)checkProxy()  {
	
}

func (this *HoneyProxy)handleSocks4Request()error  {





	return nil
}

func (this *HoneyProxy)handleHttpsRequest()error  {


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

func (this *HoneyProxy)handleConn(conn net.Conn)error  {
	defer conn.Close()
	tmpBuffer := make([]byte,0x2000)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	//读到结束标志才行
	_, err := conn.Read(tmpBuffer)
	if err != nil{
		return err
	}
	switch tmpBuffer[0] {
	case 0x4:
		return this.handleSocks4Request()
	case 0x5:
		return this.handleSocks5Request()
	case 'P':	//post
	case 'G':	//get
	case 'D':	//delete
		return this.handleHttpRequest(conn,tmpBuffer)
	case 'C':	//connect
		return this.handleHttpsRequest()
	}
	return nil
}

func (this *HoneyProxy)ListenPort(port int,maxUser int)error  {

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
