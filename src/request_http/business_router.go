package request_http

import (
	"context"
	"github.com/caoshuyu/init-project/src/controller"
	"github.com/caoshuyu/init-project/src/controller/structure"
	"github.com/caoshuyu/kit/echo_out_tools"
	"github.com/labstack/echo/v4"
)

type businessRouterFunc struct {
}

var brf businessRouterFunc
var cont controller.Controller

func (*businessRouterFunc) initProject(ectx echo.Context) (err error) {
	//获取配置参数
	req := new(structure.InitProjectInput)
	if err = ectx.Bind(req); err != nil {
		return echo_out_tools.EchoErrorData(ectx, err, 2)
	}
	var ctx context.Context
	_, err = cont.InitProject(ctx, req)
	if err != nil {
		return echo_out_tools.EchoErrorData(ectx, err, 2)
	}
	return echo_out_tools.EchoSuccessData(ectx, "{}")
}
