package id

import (
	"github.com/google/uuid"
	"strings"
)

//New 创建无分隔符的uuid v4 用于id
//
//返回string类型
func New() string {

	return uuid.New().String()
}

//NewUid 创建无分隔符的uuid v4 用于id
//
//返回uuid.UUID类型
func NewUid() uuid.UUID {

	return uuid.New()
}

//NewToken 创建无分隔符的uuid v4 用于token
//
//返回string类型
func NewToken() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

//RecoverTokenToUuid token转为 id类型
//
//返回string类型
func RecoverTokenToUuid(token string) string {
	return token[:8] + "-" + token[8:12] + "-" + token[12:16] + "-" + token[16:20] + "-" + token[20:]
}

func Parse(u string) uuid.UUID {
	return uuid.MustParse(u)
}
