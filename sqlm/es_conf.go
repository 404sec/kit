package sqlm

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type ES struct {
	Client *elasticsearch.Client
	Cfg    *ESConfig
}
type ESConfig struct {
	Addresses             []string
	Username              string
	Password              string
	Timeout               int64
	MaxIdleConnsPerHost   int
	ResponseHeaderTimeout int64
}
