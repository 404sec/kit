package wx

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/404sec/kit/cryp"
)

type WXBizDataCrypt struct {
	AppId      string
	SessionKey string
}

// ValidateWXminiUserInfo 校验微信返回的用户数据是否合法
//
//参数为获取的个人信息数据以及sessionKey 和signature ，返回是否正确
func (w *WXBizDataCrypt) ValidateWXminiUserInfo(rawData, signature string) bool {
	signature2 := cryp.Sha1(rawData + w.SessionKey)
	return signature == signature2
}

//DncryptMiniData 解密小程序数据函数
//
//参数为加密数据以及向量返回数据为解密后的数据以及错误信息
func (w *WXBizDataCrypt) DncryptMiniData(rawData, iv string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return nil, err
	}
	key_b, err := base64.StdEncoding.DecodeString(w.SessionKey)
	if err != nil {
		return nil, err
	}
	iv_b, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	dnData, err := cryp.AESCBCD(data, key_b, iv_b)
	if err != nil {
		return nil, err
	}

	return (dnData), nil
}

//获取手机号解密后的数据包内容
type PhoneBase struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}
type PhoneInfo struct {
	PhoneBase
	Watermark Watermark `json:"watermark"`
}

func (j PhoneBase) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *PhoneBase) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

//小程序数据水印
type Watermark struct {
	Appid     string      `json:"appid"`
	Timestamp interface{} `json:"timestamp"`
}

//GetPhoneNumber 获取小程序的个人手机号
//
//参数为加密数据以及向量返回数据为解密后的数据以及错误信息
func (w *WXBizDataCrypt) GetPhoneNumber(rawData, iv string) (PhoneInfo, error) {
	data, err := w.DncryptMiniData(rawData, iv)
	var phoneInfo PhoneInfo
	if err != nil {
		return phoneInfo, err
	}
	err = json.Unmarshal(data, &phoneInfo)
	if err != nil {
		return phoneInfo, err
	}
	if phoneInfo.Watermark.Appid != w.AppId {
		return phoneInfo, errors.New("数据校验失败，appId不一致")
	}
	return phoneInfo, err
}
