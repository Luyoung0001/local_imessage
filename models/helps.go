package models

import (
	"encoding/json"
	"local_imessage/utils"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// 辅助函数群
func stringInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
func removeFromSliceUsingCopy(slice []string, target string) []string {
	index := -1
	for i, s := range slice {
		if s == target {
			index = i
			break
		}
	}
	if index == -1 {
		return slice // 如果目标字符串不在切片中，直接返回原始切片
	}
	// 使用 copy 函数将后面的元素向前移动，覆盖掉目标元素
	copy(slice[index:], slice[index+1:])
	return slice[:len(slice)-1] // 删除最后一个元素
}

// IdGenerator 生成 各种Id
func IdGenerator() string {
	// 获取系统当前时间
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	// 获取一个随机数
	randomInt := rand.Intn(1000)
	randomStr := strconv.Itoa(randomInt)
	// 生成ID
	return utils.Md5Encode(formattedTime + randomStr)
}

// 从 redis 中取出的序列化后的字符串进行解析,看是那个类型
// 这里主要从数据结构的方面进行判断

func StructType(value string) string {
	var data struct {
		DataType string `json:"dataType"`
	}

	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		log.Println("Error:", err)
	}
	return data.DataType

}
