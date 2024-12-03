package aws

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

const (
	DEFAULT_REGION string = "ap-northeast-1"
)

type IS3Client interface {
	PutObject(c echo.Context, sok string, src io.Reader) (*s3.PutObjectOutput, error)
}

type s3Client struct {
	svc *s3.Client
}

func NewS3Client(cfg aws.Config) IS3Client {
	return &s3Client{
		svc: s3.NewFromConfig(cfg),
	}
}

// PutObject: S3にオブジェクトをアップロードする関数
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - sok: アップロードするファイルのS3オブジェクトキー（例: "uploads/coco.png"）
//   - src: ファイルデータ
//
// return:
//   - output: s3にアップロードした際のメタデータ
//   - error: error情報
func (cs3 *s3Client) PutObject(c echo.Context, sok string, src io.Reader) (*s3.PutObjectOutput, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Debug(configs.FetchConfigStr("aws.s3.bucket.name"))
	// s3への登録情報
	input := &s3.PutObjectInput{
		Bucket: aws.String(configs.FetchConfigStr("aws.s3.bucket.name")),
		Key:    aws.String(sok),
		Body:   src,
	}

	// オプションを取得
	optFns := getS3Options()

	// S3へのアップロード
	output, err := cs3.svc.PutObject(
		context.Background(),
		input,
		optFns...,
	)

	logger.Infof("output: %v", output)

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"画像のアップロードに失敗しました。",
			wrErrors.NewUnexpectedErrorEType(),
		)
		logger.Error(wrErr)
		return nil, wrErr
	}

	return output, nil
}

// getS3Options: S3オプションの取得
//
// args:
//   - None
//
// return:
//   - []func(*s3.Options): s3のオプション
func getS3Options() []func(*s3.Options) {
	// localの時にminioに向ける
	if configs.FetchConfigStr("ENV") == "local" {
		return []func(*s3.Options){
			func(o *s3.Options) {
				o.BaseEndpoint = aws.String("http://localhost:9000")
				o.UsePathStyle = true // パススタイルURLを使用
			},
		}
	}
	return nil
}