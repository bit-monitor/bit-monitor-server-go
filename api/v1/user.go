package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddUserRegisterRecord(c *gin.Context) {
	var err error
	var r validation.AddUserRegisterRecord
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("注册失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, entity := service.AddUserRegisterRecord(r); err != nil {
		global.WM_LOG.Error("[失败]注册", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]注册", zap.Any("entity", entity))
		response.SuccessWithData(entity, c)
	}
}

func AuditUserRegisterRecord(c *gin.Context) {
	var err error
	var r validation.AuditUserRegisterRecord
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("注册审批失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, entity := service.AuditUserRegisterRecord(r); err != nil {
		global.WM_LOG.Error("[失败]注册审批", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]注册审批", zap.Any("entity", entity))
		response.SuccessWithData(entity, c)
	}
}

func GetUserRegisterRecord(c *gin.Context) {
	var err error
	var r validation.GetUserRegisterRecord
	err = c.ShouldBindQuery(&r)
	if err != nil {
		global.WM_LOG.Error("查询注册记录失败", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		if err, data := service.GetUserRegisterRecord(r); err != nil {
			global.WM_LOG.Error("查询注册记录失败", zap.Any("err", err))
			response.FailWithError(err, c)
		} else {
			response.SuccessWithData(data, c)
		}
	}
}

func Login(c *gin.Context) {
	var err error
	var r validation.Login
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("登录失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, data := service.Login(r)
	if err != nil {
		global.WM_LOG.Error("登录失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	} else {
		response.SuccessWithData(data, c)
	}
}

func GetUser(c *gin.Context) {
	var err error
	var r validation.GetUser
	err = c.ShouldBindQuery(&r)
	if err != nil {
		global.WM_LOG.Error("查询用户列表失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, data := service.GetUser(r)
	if err != nil {
		global.WM_LOG.Error("查询用户列表失败", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		response.SuccessWithData(data, c)
	}
}

func GetUserDetail(c *gin.Context) {
	var err error
	err, userId := service.GetUserIdByContext(c)
	if err != nil {
		global.WM_LOG.Error("查询用户详情失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, user := service.GetUserDetail(userId)
	if err != nil {
		global.WM_LOG.Error("查询用户详情失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	response.SuccessWithData(user, c)
}

func GetRelatedProjectList(c *gin.Context) {
	var err error
	err, userId := service.GetUserIdByContext(c)
	if err != nil {
		global.WM_LOG.Error("查询用户详情失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, data := service.GetRelatedProjectList(userId)
	if err != nil {
		global.WM_LOG.Error("查询用户详情失败", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	response.SuccessWithData(data, c)
}
