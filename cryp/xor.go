package cryp

import (
	"encoding/hex"
	"strings"
)

//XORD Xor算法解密函数
//
//加密数据以及加密key做参数，返回加密结果以及错误信息
func XORD(data string, key string) (string, error) {
	aa, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	dataLen := len(aa)
	keyLen := len(key)
	result := ""
	for i := 0; i < dataLen; i++ {
		result += string(aa[i] ^ key[i%keyLen])
	}
	return string(result), nil
}

//XORE Xor算法加密函数
//
//解密数据以及加密key做参数，返回加密结果
func XORE(data string, key string) string {
	dataLen := len(data)
	keyLen := len(key)
	result := make([]byte, 0)

	for i := 0; i < dataLen; i++ {
		result = append(result, (data[i] ^ key[i%keyLen]))
	}
	encodedStr := hex.EncodeToString(result)

	return strings.ToUpper(encodedStr)
}
