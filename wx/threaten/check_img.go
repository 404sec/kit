package threaten

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func (d *Detect) WechatDetectImg(bt []byte, fileName string) error {
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
	destURL := "https://api.weixin.qq.com/wxa/img_sec_check?access_token=" + d.AccessToken
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
