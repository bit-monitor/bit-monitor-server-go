package global

import (
	"bit.monitor.com/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	WM_CONFIG config.Server
	WM_VP     *viper.Viper
	WM_LOG    *zap.Logger
	WM_DB     *gorm.DB
	WM_REDIS  *Redis
)
