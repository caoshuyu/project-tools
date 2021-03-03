package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/caoshuyu/init-project/src/conf"
	"github.com/caoshuyu/init-project/src/controller/structure"
	"github.com/caoshuyu/kit/filetools"
	"github.com/caoshuyu/kit/stringtools"
	"os/exec"
	"runtime"
	"strings"
)

func (*Controller) InitProject(ectx context.Context, input *structure.InitProjectInput) (out *structure.InitProjectOut, err error) {
	//检测参数
	err = _checkParams(ectx, input)
	if nil != err {
		return
	}
	buildFileList := make([]buildFile, 0)
	pf := projectFile{
		importPath: input.ProjectPath + "/" + input.ProjectName,
	}
	fList, err := pf._buildInitProject(ectx, input)
	if nil != err {
		return
	}
	buildFileList = append(buildFileList, fList...)
	buildFileList = append(buildFileList, pf._buildConf(ectx, input)...)
	buildFileList = append(buildFileList, pf._buildController(ectx, input)...)
	buildFileList = append(buildFileList, pf._buildModel(ectx, input)...)
	buildFileList = append(buildFileList, pf._buildRequestHttp(ectx, input)...)
	err = pf._buildFileOrDir(ectx, buildFileList)
	fmt.Println(err)
	return
}

func _checkParams(ectx context.Context, input *structure.InitProjectInput) error {
	if strings.EqualFold(input.ProjectName, "") {
		return errors.New("project name is null")
	}
	if strings.EqualFold(input.SavePath, "") {
		return errors.New("save path is null")
	}
	if strings.EqualFold(input.ProjectPath, "") {
		return errors.New("project path is null")
	}
	if '/' == input.SavePath[len(input.SavePath)-1] {
		input.SavePath = input.SavePath[:len(input.SavePath)-1]
	}
	if '/' == input.ProjectPath[len(input.ProjectPath)-1] {
		input.ProjectPath = input.ProjectPath[:len(input.ProjectPath)-1]
	}
	return nil
}

type buildFile struct {
	fileType int //file/dir
	filePath string
	value    string
}

const (
	fileTypeDir  = 1
	fileTypeFile = 2
)

type projectFile struct {
	projectPath             string
	resourcesPath           string
	srcPath                 string
	confPath                string
	controllerPath          string
	controllerStructurePath string
	modelPath               string
	requestHttpPath         string
	importPath              string
}

func (pf *projectFile) _buildInitProject(ectx context.Context, input *structure.InitProjectInput) ([]buildFile, error) {
	buildFileList := make([]buildFile, 0)
	filePath := input.SavePath + "/" + input.ProjectName
	//make project dir
	haveDir := false
	if filetools.CheckFileExist(filePath) {
		haveDir = true
	}
	pf.projectPath = filePath
	if haveDir {
		buildFileList = append(buildFileList, buildFile{
			fileType: fileTypeDir,
			filePath: pf.projectPath,
		})
	}
	pf.resourcesPath = filePath + "/resources"
	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeDir,
		filePath: pf.resourcesPath,
	})
	//make src dir
	pf.srcPath = filePath + "/src"
	if haveDir {
		if filetools.CheckFileExist(pf.srcPath) {
			return nil, errors.New("project dir is already used")
		}
	}
	pf.confPath = pf.srcPath + "/conf"
	pf.controllerPath = pf.srcPath + "/controller"
	pf.controllerStructurePath = pf.controllerPath + "/structure"
	pf.modelPath = pf.srcPath + "/model"
	pf.requestHttpPath = pf.srcPath + "/request_http"
	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeDir,
		filePath: pf.srcPath,
	}, buildFile{
		fileType: fileTypeDir,
		filePath: pf.confPath,
	}, buildFile{
		fileType: fileTypeDir,
		filePath: pf.controllerPath,
	}, buildFile{
		fileType: fileTypeDir,
		filePath: pf.controllerStructurePath,
	}, buildFile{
		fileType: fileTypeDir,
		filePath: pf.modelPath,
	}, buildFile{
		fileType: fileTypeDir,
		filePath: pf.requestHttpPath,
	})
	//make file
	//.gitignore
	gitignoreFile := filetools.ReadFile("./resources/gitignore.txt")
	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeFile,
		filePath: pf.projectPath + "/.gitignore",
		value:    gitignoreFile,
	})
	//go.mod
	gomodFile := make([]string, 0)
	gomodFile = append(gomodFile, "module "+input.ProjectPath+"/"+input.ProjectName, "")
	gomodFile = append(gomodFile, "go "+strings.Replace(runtime.Version(), "go", "", -1), "")
	gomodFile = append(gomodFile,
		"require (",
		"github.com/BurntSushi/toml v0.3.1",
		"github.com/caoshuyu/kit "+conf.ConfRead{}.GetKitVersion(),
		"github.com/labstack/echo/v4 v4.2.0",
		")")
	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeFile,
		filePath: pf.projectPath + "/go.mod",
		value:    strings.Join(gomodFile, "\n"),
	})
	//.toml
	tomlFile := make([]string, 0)
	switch len(input.Mysql) {
	case 0:
		//不使用mysql
	case 1:
		//单一mysql
		tomlFile = append(tomlFile,
			`# MySQL 配置`,
			`[mysql]`,
			`#用户名`,
			`username = "root"`,
			`#密码`,
			`password = "root"`,
			`#链接地址`,
			`address = "127.0.0.1:3306"`,
			`#数据库名称`,
			`db_name = "`+input.ProjectName+`"`,
			`#附加请求参数`,
			`params = "clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"`,
			`#最大连接数`,
			`max_open = 100`,
			`#最大空闲数`,
			`max_idle = 100`,
			`#连接生命时长，秒`,
			`max_lifetime = 300`,
			``,
		)
	default:
		//超过一个mysql

	}

	tomlFile = append(tomlFile, `[http]`, `port = 0`, ``)
	tomlFile = append(tomlFile, `[log]`, `save_path =""`, ``)
	tomlFile = append(tomlFile, `#配置更新授权`, `[conf_key]`, `ak = ""`, `sk = ""`, ``)
	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeFile,
		filePath: pf.projectPath + "/" + input.ProjectName + ".toml",
		value:    strings.Join(tomlFile, "\n"),
	})
	//README
	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeFile,
		filePath: pf.projectPath + "/" + "README.md",
		value:    strings.Join([]string{input.ProjectName, "==="}, "\n"),
	})
	//main.go
	mainFile := make([]string, 0)
	mainFile = append(mainFile,
		"package main",
		"",
		"import (",
		`"`+pf.importPath+`/src/conf"`,
		`"`+pf.importPath+`/src/request_http"`,
		`"`+pf.importPath+`/src/controller"`,
		`)`,
		``,
		`func main() {`,
		`//初始化配置信息`,
		`conf.InitConf()`,
		`//初始化数据库信息`,
		`controller.InitDb()`,
		`//启动HTTP服务`,
		`request_http.ListeningHTTP()`,
		`}`,
	)
	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeFile,
		filePath: pf.projectPath + "/" + "main.go",
		value:    strings.Join(mainFile, "\n"),
	})
	return buildFileList, nil
}

func (pf *projectFile) _buildConf(ectx context.Context, input *structure.InitProjectInput) []buildFile {
	buildFileList := make([]buildFile, 0)
	confileFile := make([]string, 0)
	configStructFile := make([]string, 0)
	constantFile := make([]string, 0)

	confileFile = append(confileFile,
		`package conf`,
		``,
		`import (`,
		`"errors"`,
		`"github.com/BurntSushi/toml"`,
		`"github.com/caoshuyu/kit/dlog"`,
	)
	if len(input.Mysql) > 0 {
		confileFile = append(confileFile, `"github.com/caoshuyu/kit/mysqltools"`)
	}
	confileFile = append(confileFile,
		`"strconv"`,
		`"time"`,
		`)`,
		``,
	)

	confileFile = append(confileFile, "var requestHttpPort string")
	switch len(input.Mysql) {
	case 0:
	case 1:
		confileFile = append(confileFile, "var mysql *mysqltools.MySqlConf")
	default:

	}
	confileFile = append(confileFile, "var confKey confKeyConf", "")
	confileFile = append(confileFile,
		"func InitConf() {",
		`	conf, err := getConf()`,
		`	if nil != err {`,
		`		panic(err)`,
		`	}`,
		`	requestHttpPort = ":" + strconv.Itoa(conf.Http.Port)`,
		`	initLog(&conf)`,
		`	initConfKey(&conf)`,
	)
	switch len(input.Mysql) {
	case 0:
	case 1:
		confileFile = append(confileFile, `mysql = initMysql(&conf)`)
	default:

	}
	confileFile = append(confileFile, `}`, ``)

	confileFile = append(confileFile,
		`func getConf() (config, error) {`,
		`	conf := config{}`,
		`	_, err := toml.DecodeFile("./"+SERVER_NAME+".toml", &conf)`,
		`	if nil != err {`,
		`		return conf, err`,
		`	}`,
		`	return conf, nil`,
		`}`,
		``,
	)

	confileFile = append(confileFile,
		`func initLog(conf *config) {`,
		`	dlog.SetLog(dlog.SetLogConf{`,
		`		LogType: dlog.LOG_TYPE_LOCAL,`,
		`		LogPath: conf.Log.SavePath,`,
		`		Prefix:  SERVER_NAME,`,
		`	})`,
		`}`,
		``,
	)

	switch len(input.Mysql) {
	case 0:
	case 1:
		confileFile = append(confileFile,
			`func initMysql(conf *config) *mysqltools.MySqlConf {`,
			`	return &mysqltools.MySqlConf{`,
			`		DbDsn:       conf.Mysql.Username + ":" + conf.Mysql.Password + "@tcp(" + conf.Mysql.Address + ")/" + conf.Mysql.DbName + "?" + conf.Mysql.Params,`,
			`		MaxOpen:     conf.Mysql.MaxOpen,`,
			`		MaxIdle:     conf.Mysql.MaxIdle,`,
			`		DbName:      conf.Mysql.DbName,`,
			`		MaxLifetime: time.Duration(conf.Mysql.MaxLifetime) * time.Second,`,
			`	}`,
			`}`,
			``,
		)
	default:

	}

	confileFile = append(confileFile,
		`func initConfKey(conf *config) {`,
		`	confKey = confKeyConf{`,
		`		Ak: conf.ConfKey.Ak,`,
		`		Sk: conf.ConfKey.Sk,`,
		`	}`,
		`}`,
		``,
	)

	upList := make([]string, 0)
	upList = append(upList,
		`//更新某个特定配置`,
		`func UpdateConf(confName string) (err error) {`,
		`	switch confName {`,
	)
	caseStr := make([]string, 0)
	switch len(input.Mysql) {
	case 0:
	case 1:
		caseStr = append(caseStr, `"mysql"`)
	default:

	}
	caseStr = append(caseStr, `"log"`)
	switch len(input.Redis) {
	case 0:
	case 1:
		caseStr = append(caseStr, `"redis"`)
	default:

	}
	upList = append(upList, `	case `+strings.Join(caseStr, ",")+":")
	upList = append(upList,
		`	default:`,
		`		err = errors.New("conf name not have")`,
		`		return`,
		`	}`,
		`	conf, err := getConf()`,
		`	if nil != err {`,
		`		return err`,
		`	}`,
		`	switch confName {`,
	)
	switch len(input.Mysql) {
	case 0:
	case 1:
		upList = append(upList,
			`	case "mysql":`,
			`		mysql = initMysql(&conf)`)
	default:

	}

	upList = append(upList, `	case "log":`, `		initLog(&conf)`, `	}`)

	upList = append(upList, `	return`, ``, `}`)
	confileFile = append(confileFile, upList...)

	confileFile = append(confileFile, `type ConfRead struct {}`)

	switch len(input.Mysql) {
	case 0:
	case 1:
		confileFile = append(confileFile,
			`//新配置未生效MySQL配置`,
			`func (ConfRead) NewConfGetMysqlConf() (*mysqltools.MySqlConf, error) {`,
			`	conf, err := getConf()`,
			`	if nil != err {`,
			`		return nil, err`,
			`	}`,
			`	return initMysql(&conf), nil`,
			`}`,
			``,
			`func (ConfRead) GetMysqlConf() *mysqltools.MySqlConf {`,
			`	return mysql`,
			`}`,
			``,
		)
	default:

	}

	confileFile = append(confileFile,
		`func (ConfRead) GetRequestHttpPort() string {`,
		`	return requestHttpPort`,
		`}`,
		``,
	)
	confileFile = append(confileFile,
		`func (ConfRead) GetConfKey() (ak, sk string) {`,
		`	return confKey.Ak, confKey.Sk`,
		`}`,
		``,
	)

	configStructFile = append(configStructFile, `package conf`, ``)
	configStructFile = append(configStructFile, `type config struct {`)
	switch len(input.Mysql) {
	case 0:
	case 1:
		configStructFile = append(configStructFile, `	Mysql     mysqlConf`)
	default:
	}
	configStructFile = append(configStructFile,
		`	Http      httpConf`,
		`	Log       logConf`,
		`	ConfKey   confKeyConf`,
	)
	configStructFile = append(configStructFile, `}`)

	switch len(input.Mysql) {
	case 0:
	case 1:
		configStructFile = append(configStructFile,
			`type mysqlConf struct {`,
			`	Username    string`,
			`	Password    string`,
			`	Address     string`,
			"	DbName      string `toml:\"db_name\"`",
			"	Params      string",
			"	MaxOpen     int `toml:\"max_open\"`",
			"	MaxIdle     int `toml:\"max_idle\"`",
			"	MaxLifetime int `toml:\"max_lifetime\"`",
			`}`,
			``,
		)
	default:

	}
	configStructFile = append(configStructFile,
		`type httpConf struct {`,
		`	Port int`,
		`}`,
		``,
	)
	configStructFile = append(configStructFile,
		`type logConf struct {`,
		"	SavePath string `toml:\"save_path\"`",
		"}",
		``,
	)
	configStructFile = append(configStructFile,
		`type confKeyConf struct {`,
		`	Ak string`,
		`	Sk string`,
		`}`,
		``,
	)

	constantFile = append(constantFile, "package conf", "")
	constantFile = append(constantFile, `const SERVER_NAME = "`+input.ProjectName+`"`)

	buildFileList = append(buildFileList,
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.confPath + "/config.go",
			value:    strings.Join(confileFile, "\n"),
		},
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.confPath + "/config_struct.go",
			value:    strings.Join(configStructFile, "\n"),
		},
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.confPath + "/constant.go",
			value:    strings.Join(constantFile, "\n"),
		},
	)

	return buildFileList
}

func (pf *projectFile) _buildController(ectx context.Context, input *structure.InitProjectInput) []buildFile {
	buildFileList := make([]buildFile, 0)
	baseFile := make([]string, 0)
	confServiceFile := make([]string, 0)

	baseFile = append(baseFile, `package controller`, ``)
	if len(input.Mysql) > 0 {
		baseFile = append(baseFile, `import "`+pf.importPath+`/src/model"`)
	}
	baseFile = append(baseFile, "type Controller struct {}")
	baseFile = append(baseFile, "//初始化数据库信息", "func InitDb() {")
	switch len(input.Mysql) {
	case 0:
	case 1:
		baseFile = append(baseFile, `//初始化MySQL数据库信息`, `model.Get`+stringtools.InitialUpdateStr(input.Mysql[0])+`Db()`)
	default:

	}
	baseFile = append(baseFile, "}")

	if !input.OpenUpdateConf {
		buildFileList = append(buildFileList,
			buildFile{
				fileType: fileTypeFile,
				filePath: pf.controllerPath + "/base.go",
				value:    strings.Join(baseFile, "\n"),
			},
		)
		return buildFileList
	}

	confServiceFile = append(confServiceFile, `package controller`, ``)
	confServiceFile = append(confServiceFile, `import (`)
	confServiceFile = append(confServiceFile, `"errors"`, `"strings"`, `"`+pf.importPath+`/src/conf"`)
	if len(input.Mysql) > 0 {
		confServiceFile = append(confServiceFile, `"`+pf.importPath+`/src/model"`)
	}
	confServiceFile = append(confServiceFile, `)`, ``)
	confServiceFile = append(confServiceFile, `func UpdateConf(ak, sk, name string) error {`)
	confServiceFile = append(confServiceFile, `//校验配置信息`,
		`var igcr conf.ConfRead`,
		`var err error`,
		`useAk, useSk := igcr.GetConfKey()`,
		`if !strings.EqualFold(ak, useAk) || !strings.EqualFold(sk, useSk) {`,
		`return errors.New("ak or sk error")`,
		`}`,
		`switch name {`,
	)
	confServiceFile = append(confServiceFile,
		`case "log":`,
		`err = conf.UpdateConf(name)`,
		`if nil != err {`,
		`return err`,
		`}`,
		``,
	)
	switch len(input.Mysql) {
	case 0:
	case 1:
		confServiceFile = append(confServiceFile,
			`case "mysql":`,
			`//检测新数据连接是否可用`,
			`newConf, err := conf.ConfRead{}.NewConfGetMysqlConf()`,
			`if nil != err {`,
			`return err`,
			`}`,
			`client, err := model.ConnectMysqlDb(newConf)`,
			`if nil != err {`,
			`return err`,
			`}`,
			`//更新数据库链接`,
			`err = conf.UpdateConf(name)`,
			`if nil != err {`,
			`return err`,
			`}`,
			`model.Update`+stringtools.InitialUpdateStr(input.Mysql[0])+`Db(client)`,
			``,
		)
	default:

	}

	//switch end
	confServiceFile = append(confServiceFile, "}")

	confServiceFile = append(confServiceFile, `return nil`)
	confServiceFile = append(confServiceFile, `}`)

	buildFileList = append(buildFileList,
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.controllerPath + "/base.go",
			value:    strings.Join(baseFile, "\n"),
		},
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.controllerPath + "/conf_service.go",
			value:    strings.Join(confServiceFile, "\n"),
		},
	)
	return buildFileList
}

func (pf *projectFile) _buildModel(ectx context.Context, input *structure.InitProjectInput) []buildFile {
	buildFileList := make([]buildFile, 0)

	buildFileList = append(buildFileList, pf._buildModelConnect(ectx, input)...)
	//mysql
	if len(input.Mysql) > 0 {
		buildFileList = append(buildFileList, pf._buildModelMysql(ectx, input.Mysql)...)
	}
	return buildFileList
}

func (pf *projectFile) _buildModelConnect(ectx context.Context, input *structure.InitProjectInput) []buildFile {
	buildFileList := make([]buildFile, 0)

	connectFile := make([]string, 0)
	connectFile = append(connectFile, `package model`, ``)
	connectFile = append(connectFile, `import (`)
	if len(input.Mysql) > 0 {
		connectFile = append(connectFile, `"github.com/caoshuyu/kit/mysqltools"`)
	}
	connectFile = append(connectFile, `)`, ``)

	if len(input.Mysql) > 0 {
		connectFile = append(connectFile,
			`func ConnectMysqlDb(conf *mysqltools.MySqlConf) (client *mysqltools.MysqlClient, err error) {`,
			`client = &mysqltools.MysqlClient{`,
			`Conf: conf,`,
			`}`,
			`err = client.Connect()`,
			`if nil != err {`,
			`return`,
			`}`,
			`return`,
			`}`,
		)
	}

	buildFileList = append(buildFileList,
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.modelPath + "/connect.go",
			value:    strings.Join(connectFile, "\n"),
		},
	)

	return buildFileList
}

func (pf *projectFile) _buildModelMysql(ectx context.Context, mysqlList []string) []buildFile {
	buildFileList := make([]buildFile, 0)
	mysqldbFile := make([]string, 0)
	mysqldbFile = append(mysqldbFile, `package model`, ``)
	mysqldbFile = append(mysqldbFile, `import (`)
	mysqldbFile = append(mysqldbFile,
		`"database/sql"`,
		`"`+pf.importPath+`/src/conf"`,
		`"github.com/caoshuyu/kit/mysqltools"`,
	)
	mysqldbFile = append(mysqldbFile, `)`, ``)

	//声明变量
	for _, dbName := range mysqlList {
		mysqldbFile = append(mysqldbFile, `var `+dbName+`DbClient *mysqltools.MysqlClient`)
	}
	mysqldbFile = append(mysqldbFile, ``)

	//实现方法
	for _, dbName := range mysqlList {
		mysqldbFile = append(mysqldbFile, pf._getOneMysqlDbCode(ectx, dbName)...)
	}

	buildFileList = append(buildFileList, buildFile{
		fileType: fileTypeFile,
		filePath: pf.modelPath + "/mysqldb.go",
		value:    strings.Join(mysqldbFile, "\n"),
	})
	return buildFileList
}

func (pf *projectFile) _getOneMysqlDbCode(ectx context.Context, dbName string) (funcCode []string) {

	funcCode = append(funcCode,
		`//获取链接`,
		`func Get`+stringtools.InitialUpdateStr(dbName)+`Db() *sql.DB {`,
		`	if nil == `+dbName+`DbClient {`,
		`		client, err := ConnectMysqlDb(conf.ConfRead{}.GetMysqlConf())`,
		`		if nil != err {`,
		`			panic(err)`,
		`		}`,
		`		`+dbName+`DbClient = client`,
		`	}`,
		`	return `+dbName+`DbClient.Client`,
		`}`,
	)
	funcCode = append(funcCode,
		`func Update`+stringtools.InitialUpdateStr(dbName)+`Db(newClient *mysqltools.MysqlClient) {`,
		`	`+dbName+`DbClient = newClient`,
		`}`,
	)
	return funcCode
}

func (pf *projectFile) _buildRequestHttp(ectx context.Context, input *structure.InitProjectInput) []buildFile {
	buildFileList := make([]buildFile, 0)
	httpFile := make([]string, 0)
	businessRouterFile := make([]string, 0)
	businessRouterStructFile := make([]string, 0)
	confRouterFile := make([]string, 0)

	httpFile = append(httpFile, `package request_http`, ``)
	httpFile = append(httpFile,
		`import (`,
		`	"fmt"`,
		`	"`+pf.importPath+`/src/conf"`,
		`	"github.com/caoshuyu/kit/echomiddleware"`,
		`	"github.com/labstack/echo/v4"`,
		`	"github.com/labstack/echo/v4/middleware"`,
		`	"net/http"`,
		`	"time"`,
		`)`,
	)
	httpFile = append(httpFile,
		`func ListeningHTTP() {`,
		`	e := echo.New()`,
		`	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{`,
		`		AllowOrigins: []string{`,
		`			"*",`,
		`		},`,
		`		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},`,
		`		AllowHeaders:     []string{"Accept", "Content-Type", "Authorization"},`,
		`		AllowCredentials: true,`,
		`	}))`,
		`	e.GET("/ping", func(context echo.Context) error {`,
		`		return context.JSON(http.StatusOK, "ping")`,
		`	})`,
		``,
		`	router(e)`,
		`	err := e.StartServer(&http.Server{`,
		`		Addr:              conf.ConfRead{}.GetRequestHttpPort(),`,
		`		ReadTimeout:       time.Second * 5,`,
		`		ReadHeaderTimeout: time.Second * 2,`,
		`		WriteTimeout:      time.Second * 30,`,
		`	})`,
		`	if nil != err {`,
		`		fmt.Println("server_start_err", err)`,
		`	}`,
		`}`,
		``,
	)

	httpFile = append(httpFile,
		`func router(e *echo.Echo) {`,
		`	e.Use(echomiddleware.Gls, echomiddleware.Access, echomiddleware.Recover)`,
		`	businessRouter(e)`,
	)
	if input.OpenUpdateConf {
		httpFile = append(httpFile, `	confRouter(e)`)
	}
	httpFile = append(httpFile, `}`)

	httpFile = append(httpFile,
		`//业务接口`,
		`func businessRouter(e *echo.Echo) {`,
		`}`,
	)
	if input.OpenUpdateConf {
		httpFile = append(httpFile,
			`//配置接口`,
			`func confRouter(e *echo.Echo) {`,
			`	g := e.Group("/conf")`,
			`	g.GET("/update_conf", crf.updateConf)`,
			`}`,
		)
	}

	businessRouterFile = append(businessRouterFile, `package request_http`, ``)

	businessRouterFile = append(businessRouterFile, `type businessRouterFunc struct {}`)
	businessRouterFile = append(businessRouterFile, `var brf businessRouterFunc`)

	businessRouterStructFile = append(businessRouterStructFile, `package request_http`, ``)

	buildFileList = append(buildFileList,
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.requestHttpPath + "/http.go",
			value:    strings.Join(httpFile, "\n"),
		},
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.requestHttpPath + "/business_router.go",
			value:    strings.Join(businessRouterFile, "\n"),
		},
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.requestHttpPath + "/business_router_struct.go",
			value:    strings.Join(businessRouterStructFile, "\n"),
		},
	)
	if !input.OpenUpdateConf {
		return buildFileList
	}

	confRouterFile = append(confRouterFile, `package request_http`, ``)
	confRouterFile = append(confRouterFile, `import (`)
	confRouterFile = append(confRouterFile, `	"`+pf.importPath+`/src/controller"`)
	confRouterFile = append(confRouterFile, `	"github.com/caoshuyu/kit/echo_out_tools"`)
	confRouterFile = append(confRouterFile, `	"github.com/labstack/echo/v4"`)
	confRouterFile = append(confRouterFile, `)`)
	confRouterFile = append(confRouterFile, `type confRouterFunc struct {}`)
	confRouterFile = append(confRouterFile, `var crf confRouterFunc`)
	confRouterFile = append(confRouterFile,
		`func (*confRouterFunc) updateConf(ectx echo.Context) (err error) {`,
		`	//获取校验参数`,
		`	ak := ectx.Request().Header.Get("ak")`,
		`	sk := ectx.Request().Header.Get("sk")`,
		`	name := ectx.QueryParam("name")`,
		`	err = controller.UpdateConf(ak, sk, name)`,
		`	if nil != err {`,
		`		return echo_out_tools.EchoErrorData(ectx, err, 2)`,
		`	}`,
		`	return echo_out_tools.EchoSuccessData(ectx, "")`,
		`}`,
	)

	buildFileList = append(buildFileList,
		buildFile{
			fileType: fileTypeFile,
			filePath: pf.requestHttpPath + "/conf_router.go",
			value:    strings.Join(confRouterFile, "\n"),
		},
	)

	return buildFileList
}

func (pf *projectFile) _buildFileOrDir(ectx context.Context, buildFileList []buildFile) error {
	//check path
	for _, value := range buildFileList {
		if strings.EqualFold(value.filePath, "") {
			panic("have file not have path")
		}
	}
	//make file
	for _, value := range buildFileList {
		switch value.fileType {
		case fileTypeDir:
			fmt.Println(value.fileType, value.filePath)
			err := filetools.MakeDir(value.filePath)
			if nil != err {
				return err
			}
		case fileTypeFile:
			fmt.Println(value.fileType, value.filePath)
			err := filetools.WriteFileCover(value.filePath, value.value)
			if nil != err {
				return err
			}
		}
	}

	c := exec.Command("gofmt", "-w", pf.projectPath)
	err := c.Start()
	if nil != err {
		return err
	}
	return nil
}
