package sqlm

import (
	"context"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"net"
	"net/http"
	"time"
)

type ESConfig struct {
	Addresses             []string
	Username              string
	Password              string
	Timeout               int64
	MaxIdleConnsPerHost   int
	ResponseHeaderTimeout int64
}

func (e *ESConfig) newConfig() (conf elasticsearch.Config, err error) {

	cfg := elasticsearch.Config{
		Addresses: e.Addresses, //"https://es-5e0graix.public.tencentelasticsearch.com:9200"},
		Username:  e.Username,
		Password:  e.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   e.MaxIdleConnsPerHost,
			ResponseHeaderTimeout: time.Duration(e.ResponseHeaderTimeout) * time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Duration(e.Timeout) * time.Second}).DialContext,
		},
	}

	return cfg, nil
}

//
func (e *ESConfig) NewConnect(ctx context.Context) (client *elasticsearch.Client, err error) {
	conf, _ := e.newConfig()

	client, err = elasticsearch.NewClient(conf)

	if err != nil {
		return client, err
	}
	_, err = client.Ping()

	return client, err
}
