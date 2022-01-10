package wx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/404sec/log"
	"io"
	"net/http"
	"net/url"
)

//微信的appid与Secret
type APP struct {
	Appid  string `json:"appid"`
	Secret string `json:"secret"`
}
type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (w APP) Check(code string, ctx context.Context) error {
	_, err := w.MiniLogin(code, ctx)
	return err
}

//MiniLogin 微信小程序登陆函数
//
//参数为前端获取的code 以及appid 和secretkey
//返回登陆后微信返回的 WXLoginResp 以及是否登陆错误
func (w *APP) MiniLogin(code string, ctx context.Context) (WXLoginResp, error) {
	urls := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	urls = fmt.Sprintf(urls, w.Appid, w.Secret, code)
	wxResp := WXLoginResp{}

	// 创建http get请求
	resp, err := http.Get(urls)
	if err != nil {
		return wxResp, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(ctx, "关闭http请求错误", err.Error())
		}
	}(resp.Body)

	// 解析http请求中body 数据到我们定义的结构体中
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return wxResp, err
	}

	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		return wxResp, errors.New(fmt.Sprintf("ErrCode:%s  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
	}

	return wxResp, nil
}

//GetMiniAccessToken 获取小程序access_token
//
//返回access_token 以及错误信息
func (w *APP) GetMiniAccessToken(ctx context.Context) (string, error) {
	u, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	var acc string
	if err != nil {
		return acc, err
	}
	paras := &url.Values{}
	//设置请求参数
	paras.Set("appid", w.Appid)
	paras.Set("secret", w.Secret)
	paras.Set("grant_type", "client_credential")
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	//关闭资源
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error(ctx, "关闭http请求错误", err.Error())
			}
		}(resp.Body)
	}
	if err != nil {
		return acc, errors.New("request token err :" + err.Error())
	}

	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return acc, errors.New("request token response json parse err :" + err.Error())
	}
	if jMap["errcode"] == nil || jMap["errcode"] == 0 {
		acc, _ = jMap["access_token"].(string)
		return acc, nil
	} else {
		//返回错误信息
		errcode := jMap["errcode"].(string)
		errmsg := jMap["errmsg"].(string)
		err = errors.New(errcode + ":" + errmsg)
		return acc, err
	}

}
