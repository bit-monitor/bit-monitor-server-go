package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetStatisticOverallByTimeRange(c *gin.Context) {
	var err error
	var r validation.GetStatisticOverallByTimeRange
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]获取总览页信息", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetStatisticOverallByTimeRange(r); err != nil {
		global.WM_LOG.Error("[失败]获取总览页信息", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]获取总览页信息", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetLogCountByHours(c *gin.Context) {
	var err error
	var r validation.GetLogCountByHours
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]按小时间隔获取各小时内的日志数量", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetLogCountByHours(r); err != nil {
		global.WM_LOG.Error("[失败]按小时间隔获取各小时内的日志数量", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]按小时间隔获取各小时内的日志数量", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetLogCountByDays(c *gin.Context) {
	var err error
	var r validation.GetLogCountByDays
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]按天间隔获取各日期内的日志数量", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetLogCountByDays(r); err != nil {
		global.WM_LOG.Error("[失败]按天间隔获取各日期内的日志数量", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]按天间隔获取各日期内的日志数量", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetLogCountBetweenDiffDate(c *gin.Context) {
	var err error
	var r validation.GetLogCountBetweenDiffDate
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]按天间隔获取两个日期之间的对比数据", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetLogCountBetweenDiffDate(r); err != nil {
		global.WM_LOG.Error("[失败]按天间隔获取两个日期之间的对比数据", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]按天间隔获取两个日期之间的对比数据", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetLogDistributionBetweenDiffDate(c *gin.Context) {
	var err error
	var r validation.GetLogDistributionBetweenDiffDate
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("获取两个日期之间的设备、操作系统、浏览器、网络类型的统计数据失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetLogDistributionBetweenDiffDate(r); err != nil {
		global.WM_LOG.Error("[失败]获取两个日期之间的设备、操作系统、浏览器、网络类型的统计数据", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]获取两个日期之间的设备、操作系统、浏览器、网络类型的统计数据", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetAllProjectOverviewListBetweenDiffDate(c *gin.Context) {
	var err error
	var r validation.GetAllProjectOverviewListBetweenDiffDate
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("获取用户关联的所有项目的统计情况列表失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetAllProjectOverviewListBetweenDiffDate(c, r); err != nil {
		global.WM_LOG.Error("[失败]获取用户关联的所有项目的统计情况列表", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]获取用户关联的所有项目的统计情况列表", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
