package request_http

import (
	"fmt"
	"github.com/caoshuyu/init-project/src/conf"
	"github.com/caoshuyu/kit/echomiddleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func ListeningHTTP() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{"Accept", "Content-Type", "Authorization", "Scheduler"},
		AllowCredentials: true,
	}))
	e.GET("/ping", func(context echo.Context) error {
		return context.JSON(http.StatusOK, "ping")
	})

	router(e)

	err := e.StartServer(&http.Server{
		Addr:              conf.ConfRead{}.GetRequestHttpPort(),
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      time.Second * 30,
	})
	if nil != err {
		fmt.Println("server_start_err", err)
	}

}

func router(e *echo.Echo) {
	e.Use(echomiddleware.Gls, echomiddleware.Access, echomiddleware.Recover)
	businessRouter(e)
}

//业务接口
func businessRouter(e *echo.Echo) {
	//初始化项目
	e.POST("/init_project", brf.initProject)
	//生成项目可用model文件
	e.POST("/build_model_file", brf.buildModelFile)
}
