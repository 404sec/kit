package cryp

import (
	"crypto/sha1"
	"fmt"
)

//Sha1 计算数据sha1值
//
//参数为需要计算数据返回计算的数据sha1值
func Sha1(str string) string {
	data := []byte(str)
	has := sha1.Sum(data)
	res := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return res
}
