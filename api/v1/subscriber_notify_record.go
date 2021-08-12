package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetSubscriberNotifyRecord(c *gin.Context) {
	var err error
	var r validation.GetSubscriberNotifyRecord
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("查询报警记录失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetSubscriberNotifyRecord(r); err != nil {
		global.WM_LOG.Error("查询报警记录失败", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("查询报警记录成功", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
