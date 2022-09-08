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

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	myconfig "github.com/NpoolPlatform/go-service-framework/pkg/config"
	ossconst "github.com/NpoolPlatform/go-service-framework/pkg/oss/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/secure"
)

var ErrOssClientNotInit = errors.New("oss client not init")

var (
	s3Client  *s3.Client
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

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               s3Config.EndPoint,
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := s3config.LoadDefaultConfig(context.Background(),
		s3config.WithRegion(s3Config.Region),
		s3config.WithHTTPClient(client),
		s3config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Config.AccessKey, s3Config.SecretKey, "")),
		s3config.WithEndpointResolverWithOptions(customResolver),
		s3config.WithClientLogMode(aws.LogRetries|aws.LogRequest),
	)
	if err != nil {
		return err
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.EndpointOptions.DisableHTTPS = true
	})

	return nil
}

// GetStringValueWithNameSpace not network invoke
func getS3Bucket() string {
	return _s3Config.Bucket
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

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
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
	s3out, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
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
