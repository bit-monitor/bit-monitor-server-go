package initialize

import (
	"bit.monitor.com/global"
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"os"
)

var ctx = context.Background()

func Redis() {
	redisConfig := global.WM_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		global.WM_LOG.Error("[失败]redis连接, err:", zap.Any("err", err))
		os.Exit(0)
	} else {
		global.WM_LOG.Info("[成功]redis连接:", zap.String("pong", pong))
		global.WM_REDIS = &global.Redis{
			Client:  client,
			Context: &ctx,
		}
	}
}
