package router

import (
	v1 "bit.monitor.com/api/v1"
	"github.com/gin-gonic/gin"
)

// InitProjectRouterPublic 公开路由，不需要权限校验
func InitProjectRouterPublic(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("project")
	{
		BaseRouter.GET("getByProjectIdentifier", v1.GetProjectByProjectIdentifier) // 根据项目标识查询项目
	}
}

// InitProjectRouterPrivate 私有路由，需要权限校验
func InitProjectRouterPrivate(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("project")
	{
		BaseRouter.PUT("add", v1.AddProject)              // 新增项目
		BaseRouter.GET("get", v1.GetProject)              // 条件查询项目
		BaseRouter.POST("update", v1.UpdateProject)       // 更新项目
		BaseRouter.DELETE("delete/:id", v1.DeleteProject) // 删除项目
	}
}
