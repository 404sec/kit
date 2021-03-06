package threaten

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
