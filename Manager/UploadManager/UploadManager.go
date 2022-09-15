package UploadManager

import "log"

var(
	Instance Uploader
)

type Uploader struct {
	Channel chan []byte
}

func (this *Uploader)PushHoneyData(data []byte)  {
	this.Channel <- data
}

func (this *Uploader)initUploader()error  {
	this.Channel = make(chan []byte,10000)
	return nil
}


func init()  {
	err := Instance.initUploader()
	if err != nil{
		log.Panicln(err)
	}
}