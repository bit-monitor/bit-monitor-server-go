package router

import (
	v1 "bit.monitor.com/api/v1"
	"github.com/gin-gonic/gin"
)

// InitAlarmRouterPrivate 私有路由，需要权限校验
func InitAlarmRouterPrivate(Router *gin.RouterGroup) {
	BaseRouter1 := Router.Group("alarm")
	{
		BaseRouter1.PUT("add", v1.AddAlarm)              // 新增预警
		BaseRouter1.POST("update", v1.UpdateAlarm)       // 编辑预警
		BaseRouter1.GET("get", v1.GetAlarm)              // 查询预警
		BaseRouter1.DELETE("delete/:id", v1.DeleteAlarm) // 删除预警
	}
	BaseRouter2 := Router.Group("alarmRecord")
	{
		BaseRouter2.GET("get", v1.GetAlarmRecord) // 查询预警记录
	}
}
