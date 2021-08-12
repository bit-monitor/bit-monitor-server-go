package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddHttpErrorLog(c *gin.Context) {
	var err error
	var r validation.AddHttpErrorLog
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("新增http异常日志失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.AddHttpErrorLog(r); err != nil {
		global.WM_LOG.Error("[失败]新增http异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]新增http异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetHttpErrorLog(c *gin.Context) {
	var err error
	var r validation.GetHttpErrorLog
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]条件查询http异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetHttpErrorLog(r); err != nil {
		global.WM_LOG.Error("[失败]条件查询http异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]条件查询http异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetHttpErrorLogByGroup(c *gin.Context) {
	var err error
	var r validation.GetHttpErrorLogByGroup
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]聚合查询http异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetHttpErrorLogByGroup(r); err != nil {
		global.WM_LOG.Error("[失败]聚合查询http异常日志", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]聚合查询http异常日志", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetLogCountByState(c *gin.Context) {
	var err error
	var r validation.GetLogCountByState
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]按status分类获取http日志数量", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetLogCountByState(r); err != nil {
		global.WM_LOG.Error("[失败]按status分类获取http日志数量", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]按status分类获取http日志数量", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
