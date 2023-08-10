package models

import (
	"context"
	"local_imessage/utils"
	"time"
)

func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	// 四小时就自动过期
	utils.Red.Set(ctx, key, val, timeTTL)
}
