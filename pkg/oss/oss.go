package oss

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/secure"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var ErrOssClientNotInit = errors.New("oss client not init")

var (
	s3Client  *s3.S3
	_s3Config s3Config
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

type s3Config struct {
	region    string
	endPoint  string
	accessKey string
	secretKey string
	bucket    string
}

func Init(
	regin, endpoint, ak, sk, bucket string,
) error {
	_s3Config = s3Config{
		region:    regin,
		endPoint:  endpoint,
		accessKey: ak,
		secretKey: sk,
		bucket:    bucket,
		// AccessKey: config.GetStringValueWithNameSpace(constant.ServiceName, constant.S3_PUB_KEY),
		// SecretKey: config.GetStringValueWithNameSpace(constant.ServiceName, constant.S3_PRI_KEY),
		// EndPoint:  config.GetStringValueWithNameSpace(constant.ServiceName, constant.S3_ENDPOINT),
		// Region:    config.GetStringValueWithNameSpace(constant.ServiceName, constant.S3_REGION),
		// Bucket:    config.GetStringValueWithNameSpace(constant.ServiceName, constant.S3_BUCKET),
	}

	return newS3Client(&_s3Config)
}

// GetStringValueWithNameSpace not network invoke
func getS3Bucket() string {
	return _s3Config.bucket
}

// NewS3Client main app init
func newS3Client(config *s3Config) error {
	creds := credentials.NewStaticCredentials(
		config.accessKey,
		config.secretKey,
		"",
	)
	sess, err := session.NewSession(&aws.Config{
		Credentials:          creds,
		Region:               aws.String(config.region),
		Endpoint:             aws.String(config.endPoint),
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

func PutObject(ctx context.Context, key string, body []byte) error {
	if s3Client == nil {
		return ErrOssClientNotInit
	}
	// encrypt or decode
	_out, err := secure.EncryptAES(body)
	if err != nil {
		return err
	}
	_, err = s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(getS3Bucket()),
		Key:    aws.String(key),
		Body:   bytes.NewReader(_out),
	})
	return err
}

func GetObject(ctx context.Context, key string) ([]byte, error) {
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

	// encrypt or decode
	return secure.DecryptAES(out)
}
