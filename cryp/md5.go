package cryp

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 返回MD5加密结果
//
//如果参数salt为nil则使用非盐加密
func MD5(data []byte, salt []byte) string {
	m5 := md5.New()
	m5.Write(data)
	if len(salt) > 0 {
		m5.Write(salt)
	}
	return hex.EncodeToString(m5.Sum(nil))

}
