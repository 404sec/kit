package wx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func WechatDetectImg(bt []byte, fileName, accessToken string) error {
	var bufReader bytes.Buffer
	//	"mime/multipart" 可以将上传文件封装

	mpWriter := multipart.NewWriter(&bufReader)
	//文件名无所谓
	//fileName := "detect"
	//字段名必须为media
	writer, err := mpWriter.CreateFormFile("media", fileName)
	if err != nil {

		return err
	}
	bts := append(bt, []byte("\r\n")...)
	reader := bytes.NewReader(bts)
	io.Copy(writer, reader)
	mpWriter.Close()

	client := http.DefaultClient
	destURL := "https://api.weixin.qq.com/wxa/img_sec_check?access_token=" + accessToken
	req, _ := http.NewRequest("POST", destURL, &bufReader)
	//从mpWriter中获取content-Type
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	vs := make(map[string]interface{})
	result, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(result, &vs)
	if err != nil {
		return err
	}
	//errcode 存在，且为0，返回通过
	if _, ok := vs["errcode"]; ok && vs["errcode"].(float64) == 0.0 {
		return nil
	}

	return errors.New("未知错误")

}

type Detect struct {
	AccessToken string ` json:"-"`
	Version     int    ` json:"version"`
	Openid      string ` json:"openid"`
	Scene       int    ` json:"scene"` //场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
	Content     string ` json:"content"`
	NickName    string ` json:"nickname"`
	Title       string ` json:"title"`
	Signature   string ` json:"signature"` //个性签名，该参数仅在资料类场景有效(scene=1)
}

type DetectRes struct {
	Errcode int          `json:"errcode,omitempty"`  //错误码
	Errmsg  string       `json:"errmsg,omitempty"`   //错误信息
	TraceId string       `json:"trace_id,omitempty"` //唯一请求标识，标记单次请求
	Result  ResultType   `json:"result"`             //综合结果
	Detail  []DetailType `json:"detail,omitempty"`   //详细检测结果
}
type DetailType struct {
	Strategy string `json:"strategy,omitempty"` //策略类型
	Errcode  int    `json:"errcode,omitempty"`  //错误码，仅当该值为0时，该项结果有效
	Suggest  string `json:"suggest,omitempty"`  //建议，有risky、pass、review三种值
	Label    int    `json:"label,omitempty"`    //命中标签枚举值，100 正常；10001 广告；20001 时政；20002 色情；20003 辱骂；20006 违法犯罪；20008 欺诈；20012 低俗；20013 版权；21000 其他
	Prob     int    `json:"prob,omitempty"`     //0-100，代表置信度，越高代表越有可能属于当前返回的标签（label）
	KeyWord  string `json:"keyword,omitempty"`  //命中的自定义关键词

}
type ResultType struct {
	Suggest string `json:"suggest,omitempty"` //建议，有risky、pass、review三种值
	Label   string `json:"label,omitempty"`   //命中标签枚举值，100 正常；10001 广告；20001 时政；20002 色情；20003 辱骂；20006 违法犯罪；20008 欺诈；20012 低俗；20013 版权；21000 其他
}

//WechatDetectContent   返回值这里忽略错误，只有成功与失败，也可以将错误抛给上层处理
func (d *Detect) WechatDetectContent() error {
	client := http.DefaultClient

	//请求需要json形式，用map包装contents

	bts, _ := json.Marshal(d)
	bts = append(bts, []byte("\r\n")...)
	//access_token在url中，内容在request body中
	resp, err := client.Post("https://api.weixin.qq.com/wxa/msg_sec_check?access_token="+d.AccessToken, "application/json", bytes.NewReader(bts))
	if err != nil {

		return err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//复用了map
	var res DetectRes
	err = json.Unmarshal(result, &res)

	if res.Errcode != 0 {
		return errors.New(res.Errmsg)
	}

	for _, v := range res.Detail {
		if v.Label >= 20002 {
			return errors.New("内容命中风险" + v.KeyWord)
		}
	}
	return nil

}
