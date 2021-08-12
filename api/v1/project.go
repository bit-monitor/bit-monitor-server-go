package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func AddProject(c *gin.Context) {
	var err error
	var r validation.AddProject
	// 因gin暂时还不支持从content-type: multipart/form-data中解析出array等复杂结构
	// ，因此这里暂时改为content-type: application/json的方式，跟java后台写的方案不一致
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]新增项目", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, entity := service.AddProject(r); err != nil {
		global.WM_LOG.Error("[失败]新增项目", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]新增项目", zap.Any("entity", entity))
		response.SuccessWithData(entity, c)
	}
}

func GetProject(c *gin.Context) {
	var err error
	var r validation.GetProject
	err = c.ShouldBindQuery(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]查询项目列表", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, data := service.GetProject(r)
	if err != nil {
		global.WM_LOG.Error("[失败]查询项目列表", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		response.SuccessWithData(data, c)
	}
}

func GetProjectByProjectIdentifier(c *gin.Context) {
	var err error
	var r validation.GetProjectByProjectIdentifier
	err = c.ShouldBindQuery(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]根据项目标识查询项目", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, data := service.GetProjectByProjectIdentifier(r)
	if err != nil {
		global.WM_LOG.Error("[失败]根据项目标识查询项目", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		response.SuccessWithData(data, c)
	}
}

func UpdateProject(c *gin.Context) {
	var err error

	var r validation.UpdateProject
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]更新项目", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}

	// 因gin暂时还不支持从content-type: application/x-www-form-urlencoded, 或content-type: multipart/form-data
	// 中解析出array等复杂结构，因此这里暂时改为单独解构
	var userList []uint64
	err = json.Unmarshal([]byte(c.PostForm("userList")), &userList)
	if err != nil {
		global.WM_LOG.Error("[失败]更新项目", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	r.UserList = userList
	if err, entity := service.UpdateProject(r); err != nil {
		global.WM_LOG.Error("[失败]更新项目", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]更新项目", zap.Any("entity", entity))
		response.SuccessWithData(entity, c)
	}
}

func DeleteProject(c *gin.Context) {
	var err error
	id := c.Param("id")
	if id == "" {
		global.WM_LOG.Error("[失败]删除项目，id不能为空", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	projectId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		global.WM_LOG.Error("[失败]删除项目", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, data := service.DeleteProject(projectId)
	if err != nil {
		global.WM_LOG.Error("[失败]删除项目", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		response.SuccessWithData(data, c)
	}
}
