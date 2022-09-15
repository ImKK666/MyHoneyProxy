package HoneyProxy

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"strconv"
)

type SocksRequest4 struct {
	Hostname string
	IP net.IP
	Port int
	Username string
}

func (this *HoneyProxy)handleSocks4Request(conn net.Conn,ctx *ProxyCtx,tmpBuffer []byte)error  {

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
	var HostName string
	var IP net.IP
	if socks4a {
		HostName, err = readUntilNull(socks4Req)
		if err != nil{
			return err
		}
	}else{
		IP = net.IPv4(ip[0], ip[1], ip[2], ip[3])
	}
	var realAddr string
	if len(HostName) == 0 {
		realAddr = net.JoinHostPort(IP.String(), strconv.Itoa(targetPort))
	}else{
		realAddr = net.JoinHostPort(HostName, strconv.Itoa(targetPort))
	}
	target, err := net.Dial("tcp", realAddr)
	if err != nil {
		return err
	}
	defer target.Close()

	//sock4连接成功
	var responseData [8]byte
	responseData[1] = byte(0x5a)
	copy(responseData[2:8], socks4ReqHeader[2:])
	_, err = conn.Write(responseData[:])
	if err != nil {
		return err
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
	case 'C':	//connect
		return this.handleHttpsRequest(conn,tmpBuffer)
	}
	return nil
}