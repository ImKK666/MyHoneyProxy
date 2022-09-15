package main

import (
	"MyHoneyProxy/HoneyProxy"
	"MyHoneyProxy/HttpHandler"
)

func main()  {
	proxy := HoneyProxy.NewHoneyProxy()
	proxy.SetReqHandler(HttpHandler.HandleRequest)
	proxy.SetRespHandler(HttpHandler.HandleResponse)
	proxy.StartProxy(8888,100)
}