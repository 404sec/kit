package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type Notify struct {
	Touser           string                       `json:"touser"`
	TemplateId       string                       `json:"template_id"`
	Page             string                       `json:"page"`
	MiniprogramState string                       `json:"miniprogram_state"`
	Lang             string                       `json:"lang"`
	Data             map[string]map[string]string `json:"data"`
}

func (n *Notify) SendMessage(accessToken string, pushData map[string]interface{}) error {

	url := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + accessToken

	bytesData, _ := json.Marshal(pushData)

	res, err := http.Post(url,
		"application/json;charset=utf-8", bytes.NewBuffer(bytesData))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ",
			zap.Error(err))
	}

	println(string(content))
	return nil
}
