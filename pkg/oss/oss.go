package oss

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	myconfig "github.com/NpoolPlatform/go-service-framework/pkg/config"
	ossconst "github.com/NpoolPlatform/go-service-framework/pkg/oss/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/secure"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var ErrOssClientNotInit = errors.New("oss client not init")

var (
	s3Client  *s3.S3
	_s3Config S3Config
	client    = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 20 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   20 * time.Second,
			ResponseHeaderTimeout: 20 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		},
	}
)

type S3Config struct {
	Region    string `json:"region"`
	EndPoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket,omitempty"`
}

func Init(storeType, bucketKey string) error {
	keyStore := myconfig.GetStringValueWithNameSpace(ossconst.S3NameSpace, storeType)
	s3Config := S3Config{}
	err := json.Unmarshal([]byte(keyStore), &s3Config)
	if err != nil {
		return err
	}
	namespace := myconfig.GetStringValueWithNameSpace("", myconfig.KeyHostname)
	s3Config.Bucket = myconfig.GetStringValueWithNameSpace(namespace, bucketKey)
	_s3Config = S3Config{
		Region:    s3Config.Region,
		EndPoint:  s3Config.EndPoint,
		AccessKey: s3Config.AccessKey,
		SecretKey: s3Config.SecretKey,
		Bucket:    s3Config.Bucket,
	}

	return newS3Client(&_s3Config)
}

// GetStringValueWithNameSpace not network invoke
func getS3Bucket() string {
	return _s3Config.Bucket
}

// NewS3Client main app init
func newS3Client(config *S3Config) error {
	creds := credentials.NewStaticCredentials(
		config.AccessKey,
		config.SecretKey,
		"",
	)
	sess, err := session.NewSession(&aws.Config{
		Credentials:          creds,
		Region:               aws.String(config.Region),
		Endpoint:             aws.String(config.EndPoint),
		DisableSSL:           aws.Bool(true),
		HTTPClient:           client,
		S3ForcePathStyle:     aws.Bool(true),
		S3Disable100Continue: aws.Bool(true),
	})
	if err != nil {
		return err
	}
	s3Client = s3.New(sess)
	return nil
}

func PutObject(ctx context.Context, key string, body []byte, encrypt bool) error {
	if s3Client == nil {
		return ErrOssClientNotInit
	}
	// encrypt or not
	if encrypt {
		_out, err := secure.EncryptAES(body)
		if err != nil {
			return err
		}
		body = _out
	}

	_, err := s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(getS3Bucket()),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	return err
}

func GetObject(ctx context.Context, key string, decrypt bool) ([]byte, error) {
	if s3Client == nil {
		return nil, ErrOssClientNotInit
	}
	s3out, err := s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(getS3Bucket()),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	defer s3out.Body.Close()

	out, err := ioutil.ReadAll(s3out.Body)
	if err != nil {
		return nil, err
	}

	// decrypt or not
	if decrypt {
		return secure.DecryptAES(out)
	}
	return out, nil
}
