package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 将 data 进行 md5 加密成一个 16 进制的串,之后将其转成字符串
// 小写

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写

func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

//随机数加密

func MakePassword(plainPWD, salt string) string {
	return Md5Encode(plainPWD + salt)

}

// 查询密码是否正确

func ValidPassword(plainPWD, salt, password string) bool {
	// password 为数据库中的密码
	// salt 为 UserBasic 中的字段
	return Md5Encode(plainPWD+salt) == password
}
