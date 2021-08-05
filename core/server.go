package core

import (
	"bit.monitor.com/global"
	"bit.monitor.com/initialize"
	"fmt"
	"go.uber.org/zap"
)

func RunWindowsServer() {
	Router := initialize.Router()
	address := fmt.Sprintf(":%d", global.WM_CONFIG.System.Addr)
	s := initServer(address, Router)
	global.WM_LOG.Info("web服务成功运行在端口号", zap.String("address", address))
	global.WM_LOG.Error(s.ListenAndServe().Error())
}
