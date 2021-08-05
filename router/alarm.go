package router

import (
	v1 "bit.monitor.com/api/v1"
	"github.com/gin-gonic/gin"
)

// InitAlarmRouterPrivate 私有路由，需要权限校验
func InitAlarmRouterPrivate(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("alarm")
	{
		BaseRouter.PUT("add", v1.AddAlarm)              // 新增预警
		BaseRouter.POST("update", v1.UpdateAlarm)       // 编辑预警
		BaseRouter.GET("get", v1.GetAlarm)              // 查询预警
		BaseRouter.DELETE("delete/:id", v1.DeleteAlarm) // 删除预警
	}
}
