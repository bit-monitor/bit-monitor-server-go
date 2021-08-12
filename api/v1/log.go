package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddLog(c *gin.Context) {
	var err error
	var data interface{}
	var r validation.AddLog
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("通用日志打点上传失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	switch r.LogType {
	case "JS_ERROR":
		var rJs validation.AddJsErrorLog
		err = c.ShouldBind(&rJs)
		if err == nil {
			err, data = service.AddJsErrorLog(rJs)
		}
		break
	case "HTTP_ERROR":
		var rHttp validation.AddHttpErrorLog
		err = c.ShouldBind(&rHttp)
		if err == nil {
			err, data = service.AddHttpErrorLog(rHttp)
		}
		break
	case "RESOURCE_LOAD_ERROR":
		var rRes validation.AddResourceLoadErrorLog
		err = c.ShouldBind(&rRes)
		if err == nil {
			err, data = service.AddResourceLoadErrorLog(rRes)
		}
		break
	case "CUSTOM_ERROR":
		var rCus validation.AddCustomErrorLog
		err = c.ShouldBind(&rCus)
		if err == nil {
			err, data = service.AddCustomErrorLog(rCus)
		}
		break
	default:
		break
	}
	if err != nil {
		global.WM_LOG.Error("[失败]通用日志打点上传", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]通用日志打点上传", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func AddClient(c *gin.Context) {
	var err error
	var r validation.AddClient
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("新增日志客户端用户失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.AddClient(r); err != nil {
		global.WM_LOG.Error("[失败]新增日志客户端用户", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]新增日志客户端用户", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func ListLog(c *gin.Context) {
	var err error
	var r validation.ListLog

	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("多条件高级查询失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}

	// 因gin暂时还不支持从content-type: application/x-www-form-urlencoded, 或content-type: multipart/form-data
	// 中解析出array等复杂结构，因此这里暂时改为单独解构
	var conditionList []validation.ConditionListItem
	err = json.Unmarshal([]byte(c.PostForm("conditionList")), &conditionList)
	if err != nil {
		global.WM_LOG.Error("多条件高级查询失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	r.ConditionList = conditionList

	if err, data := service.ListLog(r); err != nil {
		global.WM_LOG.Error("[失败]多条件高级查询", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		global.WM_LOG.Info("[成功]多条件高级查询", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
