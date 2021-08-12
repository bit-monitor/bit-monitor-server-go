package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddResourceLoadErrorLog(c *gin.Context) {
	var err error
	var r validation.AddResourceLoadErrorLog
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]新增resourceLoad异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.AddResourceLoadErrorLog(r); err != nil {
		global.WM_LOG.Error("[失败]新增resourceLoad异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]新增resourceLoad异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetResourceLoadErrorLog(c *gin.Context) {
	var err error
	var r validation.GetResourceLoadErrorLog
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]条件查询resourceLoad异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetResourceLoadErrorLog(r); err != nil {
		global.WM_LOG.Error("[失败]条件查询resourceLoad异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]条件查询resourceLoad异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetResourceLoadErrorLogByGroup(c *gin.Context) {
	var err error
	var r validation.GetResourceLoadErrorLogByGroup
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]聚合查询resourceLoad异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetResourceLoadErrorLogByGroup(r); err != nil {
		global.WM_LOG.Error("[失败]聚合查询resourceLoad异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]聚合查询resourceLoad异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetOverallByTimeRange(c *gin.Context) {
	var err error
	var r validation.GetOverallByTimeRange
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]获取总览统计信息", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetOverallByTimeRange(r); err != nil {
		global.WM_LOG.Error("[失败]获取总览统计信息", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]获取总览统计信息", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
