package main

import (
	"github.com/caoshuyu/init-project/src/conf"
	"github.com/caoshuyu/init-project/src/request_http"
)

func main() {
	//初始化配置信息
	conf.InitConf()

	//启动HTTP服务
	request_http.ListeningHTTP()
}
