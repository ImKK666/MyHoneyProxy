package Model

type HoneyReq struct {
	//GET,POST
	Method         string              `json:"method"`
	//攻击者地址,不带端口
	RemoteAddress  string              `json:"remote_address"`
	//http,https
	Scheme         string              `json:"scheme"`
	//主机网址
	Host           string              `json:"host"`
	//端口
	Port           string              `json:"port"`
	//路径
	Path           string              `json:"path"`
	//完整的Url
	Url            string              `json:"url"`
	//请求参数
	RequestParam   map[string][]string `json:"request_param"`
	//请求头
	RequestHeader  map[string][]string `json:"request_header"`
	//请求体
	RequestBody    []byte          		`json:"request_body"`
	//响应头
	ResponseHeader map[string][]string `json:"response_header"`
	//响应体
	ResponseBody   []byte              `json:"response_body"`
	//自身IP地址
	Ip             string              `json:"ip"`
	Origin         string              `json:"origin"`
	OriginDetail   string              `json:"origin_detail"`
	//代理用户名
	UserName       string			    `json:"username"`
	//代理密码
	Password       string              `json:"password"`
	//响应码
	StatusCode     int                 `json:"status_code"`
	//请求时间
	AttackTime     int64               `json:"attack_time"`
}
