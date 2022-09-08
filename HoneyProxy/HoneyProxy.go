package HoneyProxy

import (
	"bufio"
	"net"
	"strconv"
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

func (this *HoneyProxy)handleSocks4Request(bufConn *bufio.Reader)error  {





	return nil
}



func (this *HoneyProxy)handleSocks5Request(bufConn *bufio.Reader)error  {




	return nil
}

func (this *HoneyProxy)handleConn(conn net.Conn)error  {

	defer conn.Close()

	bufConn := bufio.NewReader(conn)
	version := []byte{0}
	if _, err := bufConn.Read(version); err != nil {
		return err
	}

	//socks4
	if version[0] == 0x4{
		return this.handleSocks4Request(bufConn)
	}else if version[0] == 0x5{
		return this.handleSocks5Request(bufConn)
	}


	return nil
}

func (this *HoneyProxy)ListenPort(port int)error  {
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
