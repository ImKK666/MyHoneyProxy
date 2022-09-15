package HoneyProxy

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
)

type socks5Header struct {
	version byte
	method byte
}

func socks5_readMethods(r io.Reader) ([]byte, error) {
	header := []byte{0}
	if _, err := r.Read(header); err != nil {
		return nil, err
	}
	numMethods := int(header[0])
	methods := make([]byte, numMethods)
	_, err := io.ReadAtLeast(r, methods, numMethods)
	return methods, err
}

func (this *HoneyProxy)parseSocks5Header(bufConn *bufio.Reader)(retHeader socks5Header,retErr error)  {
	var err error

	//检查版本
	retHeader.version,err = bufConn.ReadByte()
	if err != nil{
		return retHeader,err
	}
	if retHeader.version != 0x5 {
		return retHeader,errors.New("Unsupported SOCKS version")
	}

	//读取方法
	retHeader.method,err = bufConn.ReadByte()
	if err != nil{
		return retHeader,err
	}



	return retHeader,nil
}

func (this *HoneyProxy)handleSocks5Request(conn net.Conn,tmpBuffer []byte,ctx *ProxyCtx)error  {

	socks5Req := bufio.NewReader(bytes.NewReader(tmpBuffer))
	this.parseSocks5Header(socks5Req)

	return nil
}