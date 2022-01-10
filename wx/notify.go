package wx

/*
import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type WXPayNotifyReq struct {
	Appid          string  `xml:"appid"`
	Bank_type      string  `xml:"bank_type"`
	Cash_fee       float64 `xml:"cash_fee"`
	Fee_type       string  `xml:"fee_type"`
	Is_subscribe   string  `xml:"is_subscribe"`
	Mch_id         string  `xml:"mch_id"`
	Nonce_str      string  `xml:"nonce_str"`
	Openid         string  `xml:"openid"`
	Out_trade_no   string  `xml:"out_trade_no"`
	Result_code    string  `xml:"result_code"`
	Return_code    string  `xml:"return_code"`
	Sign           string  `xml:"sign"`
	Time_end       string  `xml:"time_end"`
	Total_fee      float64 `xml:"total_fee"`
	Trade_type     string  `xml:"trade_type"`
	Transaction_id string  `xml:"transaction_id"`
}

type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

func WeixinNoticeHandler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("读取http body失败，原因!", err)
		http.Error(rw.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	logger.Info("微信支付异步通知，HTTP Body:", string(body))

	var mr WXPayNotifyReq
	err = xml.Unmarshal(body, &mr)
	if err != nil {
		logger.Error("解析HTTP Body格式到xml失败，原因!", err)
		http.Error(rw.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["appid"] = mr.Appid
	reqMap["bank_type"] = mr.Bank_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce_str
	reqMap["openid"] = mr.Openid
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["result_code"] = mr.Result_code
	reqMap["return_code"] = mr.Return_code
	reqMap["time_end"] = mr.Time_end
	reqMap["total_fee"] = mr.Total_fee
	reqMap["trade_type"] = mr.Trade_type
	reqMap["transaction_id"] = mr.Transaction_id

	var resp WXPayNotifyResp
	//进行签名校验
	if wxpayVerifySign(reqMap, mr.Sign) {
		//transactionId := reqMap["transaction_id"]
		orderCode := reqMap["out_trade_no"]
		total_fee := reqMap["total_fee"].(float64) //分->元 除以100
		rows, err := mysqlDB.Query("SELECT * FROM canyin_order WHERE dno = ?", orderCode)
		if err != nil {
			logger.Error("微信查询价格错误", err)
			return
		}
		defer rows.Close()
		orders := RowResult(rows)
		if len(orders) > 0 {
			orderInfo := orders[0].(map[string]interface{})
			//orderId := ToStr(orderInfo["id"])
			allcost, _ := strconv.ParseFloat(ToStr(orderInfo["allcost"]), 64)
			logger.Info("价格比对", "---", allcost, "---", total_fee)
			//商户系统对于支付结果通知的内容一定要做签名验证,并校验返回的订单金额是否与商户侧的订单金额一致，防止数据泄漏导致出现“假通知”，造成资金损失
			if allcost == total_fee {
				logger.Info("订单验证成功")
				//以下是业务处理
			}
			resp.Return_code = "SUCCESS"
			resp.Return_msg = "OK"
		} else {
			resp.Return_code = "FAIL"
			resp.Return_msg = "无此订单"
		}
	} else {
		resp.Return_code = "FAIL"
		resp.Return_msg = "failed to verify sign, please retry!"
	}

	//结果返回，微信要求如果成功需要返回return_code "SUCCESS"
	bytes, _err := xml.Marshal(resp) //string(bytes)
	strResp := strings.Replace(bytes2str(bytes), "WXPayNotifyResp", "xml", -1)
	if _err != nil {
		logger.Error("xml编码失败，原因：", _err)
		http.Error(rw.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	rw.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(rw.(http.ResponseWriter), strResp)
}

//微信支付签名验证函数
func wxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {
	//方法名 行数
	pc, _, line, _ := runtime.Caller(0)
	fc := runtime.FuncForPC(pc)

	WECHAT_API_KEY := "" //微信商户key
	signCalc := wxpayCalcSign(needVerifyM, WECHAT_API_KEY)
	logger.Info(fc.Name(), line, "计算出来的sign: ", signCalc)
	logger.Info(fc.Name(), line, "微信异步通知sign: ", sign)
	if sign == signCalc {
		logger.Info(fc.Name(), line, "签名校验通过!")
		return true
	}

	logger.Error(fc.Name(), line, "签名校验失败!")
	return false
}

//微信支付计算签名的函数
func wxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	//方法名 行数
	pc, _, line, _ := runtime.Caller(0)
	fc := runtime.FuncForPC(pc)

	logger.Info(fc.Name(), line, "微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		logger.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write(str2bytes(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

*/
