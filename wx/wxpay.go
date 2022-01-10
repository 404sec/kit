package wx

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/404sec/log"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"go.uber.org/zap"
	"strings"
)

type Pay struct {
	Mchid                      string
	Mchcertificateserialnumber string
	Mchapiv3key                string
	Privatekeystr              string
}

func (p *Pay) Init() (client *core.Client, wxCtx context.Context, err error) {
	var (
		// 商户号
		mchID string = p.Mchid
		// 商户证书序列号
		mchCertificateSerialNumber string = p.Mchcertificateserialnumber
		// 商户APIv3密钥
		mchAPIv3Key string = p.Mchapiv3key
	)

	privateKeyStr := p.Privatekeystr
	privateKeyStr = strings.ReplaceAll(privateKeyStr, "@@@", "\n")

	mchPrivateKey, err := p.LoadPrivateKeyOwn()

	wxCtx = context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err = core.NewClient(wxCtx, opts...)
	if err != nil {
		log.Error(wxCtx, "new wechat pay client err ", zap.Error(err))
	}

	return client, wxCtx, err

}

func (p *Pay) LoadPrivateKeyOwn() (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(p.Privatekeystr))
	if block == nil {
		return nil, fmt.Errorf("decode private key err")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("the kind of PEM should be PRVATE KEY")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key err:%s", (err))
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s is not rsa private key", p.Privatekeystr)
	}
	return privateKey, nil
}
