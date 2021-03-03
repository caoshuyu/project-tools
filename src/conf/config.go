package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/caoshuyu/kit/dlog"
	"strconv"
)

var requestHttpPort string
var kitVersion string

func InitConf() {
	conf, err := getConf()
	if nil != err {
		panic(err)
	}
	requestHttpPort = ":" + strconv.Itoa(conf.Http.Port)
	initLog(&conf)
	kitVersion = conf.KitVersion
}

func getConf() (config, error) {
	conf := config{}
	_, err := toml.DecodeFile("./init-project.toml", &conf)
	if nil != err {
		return conf, err
	}
	return conf, nil
}

func initLog(conf *config) {
	dlog.SetLog(dlog.SetLogConf{
		LogType: dlog.LOG_TYPE_LOCAL,
		LogPath: conf.Log.SavePath,
		Prefix:  SERVER_NAME,
	})
}

type ConfRead struct {
}

func (ConfRead) GetRequestHttpPort() string {
	return requestHttpPort
}

func (ConfRead) GetKitVersion() string {
	return kitVersion
}
