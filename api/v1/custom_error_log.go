package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddCustomErrorLog(c *gin.Context) {
	var err error
	var r validation.AddCustomErrorLog
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]新增custom异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.AddCustomErrorLog(r); err != nil {
		global.WM_LOG.Error("[失败]新增custom异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]新增custom异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetCustomErrorLog(c *gin.Context) {
	var err error
	var r validation.GetCustomErrorLog
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]条件查询custom异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetCustomErrorLog(r); err != nil {
		global.WM_LOG.Error("[失败]条件查询custom异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]条件查询custom异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetCustomErrorLogByGroup(c *gin.Context) {
	var err error
	var r validation.GetCustomErrorLogByGroup
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]聚合查询custom异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetCustomErrorLogByGroup(r); err != nil {
		global.WM_LOG.Error("[失败]聚合查询custom异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]聚合查询custom异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
