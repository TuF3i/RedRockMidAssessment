package redis

import (
	"RedRockMidAssessment/core"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func ConnectToRedis(debug bool) error {
	//建立连接
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%v:%v",
			core.Config.Db.Redis.Addr,
			core.Config.Db.Redis.Port,
		),
		Password: core.Config.Db.Redis.Passwd,
		DB:       core.Config.Db.Redis.DefaultDB,
	})

	//测试连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	//全局化
	core.RedisConn = client

	return nil
}
