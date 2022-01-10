package threaten

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

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
