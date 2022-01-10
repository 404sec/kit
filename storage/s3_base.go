package storage

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3 struct {
	AccessKey        string
	SecretKey        string
	Endpoint         string
	Region           string
	DisableSSL       bool
	S3ForcePathStyle bool
	MaxSize          int64
}

func (s *S3) newConfig() *aws.Config {

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(s.AccessKey, s.SecretKey, ""),
		Endpoint:         aws.String(s.Endpoint),
		Region:           aws.String(s.Region),
		DisableSSL:       aws.Bool(s.DisableSSL),
		S3ForcePathStyle: aws.Bool(s.S3ForcePathStyle), //virtual-host style方式，不要修改
	}

	return s3Config
}

func (s *S3) NewConnect(ctx context.Context) (*session.Session, error) {
	conf := s.newConfig()
	client, err := session.NewSession(conf) // token can be left blank for now

	return client, err
}
