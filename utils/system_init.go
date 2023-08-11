package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("/Users/luliang/GoLand/local_imessage/config") //带绝对路径
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("config.app"))

}

// 初始化 Redis

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.minIdleConn"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	fmt.Println("config Redis:", viper.Get("redis"))
}

const (
	PublishKey = "websocket"
)

// Publish() 将消息发送到 Redis

func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish...:", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅 Redis 消息

func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe sub...:", sub)
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("上面那个函数怎么卡住了?")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("---------------------------------")
	fmt.Println("Subscribe msg PayLoad...:", msg.Payload)
	return msg.Payload, err

}
