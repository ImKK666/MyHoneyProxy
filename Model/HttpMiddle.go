package Model


type HttpMiddle struct {
	RequestBody  []byte
	ResponseBody []byte
	//攻击者IP
	IpAddr string
	//请求的域名
	Host string
	//请求的端口
	Port string
	//过滤上传
	SkipBody bool
}