package router

import (
	v1 "bit.monitor.com/api/v1"
	"github.com/gin-gonic/gin"
)

// InitSubscriberNotifyRecordRouterPrivate 私有路由，需要权限校验
func InitSubscriberNotifyRecordRouterPrivate(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("subscriberNotifyRecord")
	{
		BaseRouter.GET("get", v1.GetSubscriberNotifyRecord)                               // 查询报警记录
		BaseRouter.GET("getWithRelatedInfo", v1.GetWithRelatedInfoSubscriberNotifyRecord) // 查询报警记录-带关联信息
	}
}
