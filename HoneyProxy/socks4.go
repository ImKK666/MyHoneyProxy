package HoneyProxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

type SocksRequest4 struct {
	Hostname string
	IP net.IP
	Port int
	Username string
}

func (this *HoneyProxy)handleSocks4Request(conn net.Conn,tmpBuffer []byte,ctx *ProxyCtx)error  {

	socks4Req := bufio.NewReader(bytes.NewReader(tmpBuffer))
	var socks4ReqHeader [8]byte
	_, err := io.ReadFull(socks4Req, socks4ReqHeader[:])
	if err != nil {
		return err
	}
	//客户端访问端口
	targetPort := int(binary.BigEndian.Uint16(socks4ReqHeader[2:4]))
	ip := socks4ReqHeader[4:8]
	socks4a := ip[0] == 0 && ip[1] == 0 && ip[2] == 0 && ip[3] != 0
	ctx.ProxyAuth.UserName, err = readUntilNull(socks4Req)
	if err != nil{
		return err
	}
	var realAddr string
	var hostName string
	var IP net.IP
	if socks4a {
		hostName, err = readUntilNull(socks4Req)
		if err != nil{
			return err
		}
		realAddr = net.JoinHostPort(hostName, strconv.Itoa(targetPort))
	}else{
		IP = net.IPv4(ip[0], ip[1], ip[2], ip[3])
		realAddr = net.JoinHostPort(IP.String(), strconv.Itoa(targetPort))
	}
	log.Println(realAddr)

	//sock4连接成功
	var responseData [8]byte
	responseData[1] = byte(0x5a)
	copy(responseData[2:8], socks4ReqHeader[2:])
	_, err = conn.Write(responseData[:])
	if err != nil {
		return err
	}

	//确定是https
	if targetPort == 443{
		tlsConfig, err := TLSConfigFromCA()(hostName,ctx)
		if err != nil{
			return err
		}
		rawClientTls := tls.Server(conn,tlsConfig)
		defer rawClientTls.Close()
		err = rawClientTls.Handshake()
		if err != nil {
			return err
		}
		cReq,err := http.ReadRequest(bufio.NewReader(rawClientTls))
		if err != nil{
			return err
		}
		cReq.RemoteAddr = conn.RemoteAddr().String()
		if !httpsRegexp.MatchString(cReq.URL.String()) {
			cReq.URL, _ = url.Parse("https://" + cReq.Host + cReq.URL.String())
		}
		ctx.Req = cReq
		return this.executeHttpsRequest(rawClientTls,ctx)
	}

	//开始解析请求
	tmpBuffer,err = this.readCompleteReq(conn)
	if err != nil{
		return err
	}
	switch tmpBuffer[0] {
	case 'O':	//options
		return this.handleHttpRequest(conn,tmpBuffer,ctx)
	case 'P':	//post,put,patch
		return this.handleHttpRequest(conn,tmpBuffer,ctx)
	case 'T':	//trace
		return this.handleHttpRequest(conn,tmpBuffer,ctx)
	case 'D':	//delete
		return this.handleHttpRequest(conn,tmpBuffer,ctx)
	case 'H':	//head
		return this.handleHttpRequest(conn,tmpBuffer,ctx)
	case 'G':	//get
		return this.handleHttpRequest(conn,tmpBuffer,ctx)
	case 0x16:
		//return this.handleHttpsRequest(conn,ctx)
	case 'C':	//connect
		//return this.handleHttpsRequest(conn,ctx)
	}
	return nil
}