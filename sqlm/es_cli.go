package sqlm

import (
	"context"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"net"
	"net/http"
	"time"
)

func (e *ES) newConfig() (conf elasticsearch.Config, err error) {

	cfg := elasticsearch.Config{
		Addresses: e.Cfg.Addresses, //"https://es-5e0graix.public.tencentelasticsearch.com:9200"},
		Username:  e.Cfg.Username,
		Password:  e.Cfg.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   e.Cfg.MaxIdleConnsPerHost,
			ResponseHeaderTimeout: time.Duration(e.Cfg.ResponseHeaderTimeout) * time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Duration(e.Cfg.Timeout) * time.Second}).DialContext,
		},
	}

	return cfg, nil
}

//
func (e *ES) NewConnect(ctx context.Context) (err error) {
	conf, _ := e.newConfig()

	e.Client, err = elasticsearch.NewClient(conf)

	if err != nil {
		return err
	}
	_, err = e.Client.Ping()

	return err
}

func (e *ES) Close() error {
	return nil
}

//
func (e *ES) GetClient() *elasticsearch.Client {
	return e.Client
}
