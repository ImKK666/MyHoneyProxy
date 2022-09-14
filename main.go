package main

import (
	"MyHoneyProxy/HoneyProxy"
)

func main()  {
	proxy := HoneyProxy.NewHoneyProxy()
	proxy.SetReqHandler()
	proxy.SetRespHandler()
	proxy.ListenPort(8888,100)
}