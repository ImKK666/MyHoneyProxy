package main

import (
	"MyHoneyProxy/HoneyProxy"
	"MyHoneyProxy/HttpHandler"
	"log"
)

func main()  {
	honeyProxy := HoneyProxy.NewHoneyProxy()
	honeyProxy.SetReqHandler(HttpHandler.HandleRequest)
	honeyProxy.SetRespHandler(HttpHandler.HandleResponse)
	err := honeyProxy.StartProxy(9999)
	if err != nil{
		log.Println(err)
	}
}
