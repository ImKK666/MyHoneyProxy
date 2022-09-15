package HttpHandler

import (
	"MyHoneyProxy/HoneyProxy"
	"MyHoneyProxy/Model"
	"MyHoneyProxy/Utils"
	"bufio"
	"bytes"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const strRegexReq = `\.(css|dll|gif|ico|jpeg|jpg|js|mkv|mp4|png|svg|zip)$`

var(
	regex_Req = regexp.MustCompile(strRegexReq)
	gFilterDomain = make(map[string]struct{})
)

func parseHostPort(host string)(retHost string,retPort string)  {
	sliceHost := strings.Split(host, ":")
	if len(sliceHost) > 1 {
		retHost, retPort = sliceHost[0], sliceHost[1]
	} else {
		retHost = sliceHost[0]
		retPort = "443"
	}
	return retHost,retPort
}

//过滤非法域名,返回true表示过滤

func FilterBlackDomain(host string)bool  {
	var domain string
	var err error
	if Utils.IsIpv4Addr(host) == false{
		domain,err = publicsuffix.Domain(host)
		if err != nil{
			domain = host
		}
	}else{
		domain = host
	}
	//过滤非法域名
	_,bExists := gFilterDomain[domain]
	return bExists
}

func copyRequestBody(res *http.Request) ([]byte, error) {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return buf, nil
}

func HandleRequest(req *http.Request, ctx *HoneyProxy.ProxyCtx) (*http.Request, *http.Response) {
	httpMiddle := Model.HttpMiddle{}
	ctx.UserData = &httpMiddle
	httpMiddle.IpAddr, _, _ = net.SplitHostPort(req.RemoteAddr)
	if httpMiddle.IpAddr == ""{
		httpMiddle.IpAddr = "127.0.0.1"
	}
	//第二步,解析出Host
	httpMiddle.Host,httpMiddle.Port = parseHostPort(req.Host)
	if FilterBlackDomain(httpMiddle.Host) == true{
		httpMiddle.SkipBody = true
		return req,&http.Response{Status:"404 Forbidden", StatusCode:404, Header:make(http.Header),
			Body: ioutil.NopCloser(strings.NewReader("404 Forbidden"))}
	}
	if regex_Req.MatchString(req.URL.Path) == true{
		httpMiddle.SkipBody = true
		return req, nil
	}
	var err error
	//获取请求内容
	httpMiddle.RequestBody, err = copyRequestBody(req)
	if err != nil {
		return req, nil
	}
	return req, nil
}


func init()  {
	hBlackFile,err := os.Open("./blackList.txt")
	if err != nil{
		return
	}
	defer hBlackFile.Close()
	hScanner := bufio.NewScanner(hBlackFile)
	for hScanner.Scan(){
		gFilterDomain[hScanner.Text()]= struct{}{}
	}
	log.Println("加载黑名单域名完成:",len(gFilterDomain))
}