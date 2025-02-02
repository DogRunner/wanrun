package aws

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

const (
	DEFAULT_REGION string = "ap-northeast-1"
)

type IS3Provider interface {
	PutObject(c echo.Context, sok string, src io.Reader) error
	DeleteObject(c echo.Context, sok string) error
	GetObject(c echo.Context, sok string) error
}

type s3Provider struct {
	svc *s3.Client
}

func NewS3Provider(cfg aws.Config) IS3Provider {
	return &s3Provider{
		svc: s3.NewFromConfig(cfg),
	}
}

// PutObject: S3にオブジェクトをアップロードする関数
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: アップロードするファイルのS3オブジェクトキー（例: "uploads/coco.png"）
//   - io.Reader: ファイルデータ
//
// return:
//   - error: error情報
func (cs3 *s3Provider) PutObject(
	c echo.Context,
	sok string,
	src io.Reader,
) error {
	logger := log.GetLogger(c).Sugar()

	logger.Debugf("Bucket: %v, Key: %v", configs.FetchConfigStr("aws.s3.bucket.name"), sok)

	// s3への登録情報
	input := &s3.PutObjectInput{
		Bucket: aws.String(configs.FetchConfigStr("aws.s3.bucket.name")),
		Key:    aws.String(sok),
		Body:   src,
	}

	logger.Infof("PutObject input: %+v", input)

	// オプションを取得
	optFns := getS3Options()

	// S3へのアップロード
	output, err := cs3.svc.PutObject(
		context.Background(),
		input,
		optFns...,
	)

	logger.Infof("PutObject output: %+v", output)

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"画像のアップロードに失敗しました。",
			wrErrors.NewUnexpectedErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	return nil
}

// DeleteObject: S3にオブジェクトを削除する関数
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: アップロードするファイルのS3オブジェクトキー（例: "uploads/coco.png"）
//
// return:
//   - error: error情報
func (cs3 *s3Provider) DeleteObject(c echo.Context, sok string) error {
	logger := log.GetLogger(c).Sugar()

	logger.Debugf("Bucket: %v, Key: %v", configs.FetchConfigStr("aws.s3.bucket.name"), sok)

	// s3への削除情報
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(configs.FetchConfigStr("aws.s3.bucket.name")),
		Key:    aws.String(sok),
	}

	logger.Infof("Delete input: %+v", deleteInput)

	// オプションを取得
	optFns := getS3Options()

	// s3の対象のオブジェクト削除
	deleteOutput, err := cs3.svc.DeleteObject(
		context.Background(),
		deleteInput,
		optFns...,
	)

	logger.Infof("Delete output: %+v", deleteOutput)

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"画像、ファイルの削除に失敗しました。",
			wrErrors.NewCmsServerErrorEType(),
		)

		logger.Errorf("S3 upload failure: %v", wrErr)
		return wrErr
	}

	return nil
}

// GetObject: S3にオブジェクト検索する関数
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: アップロードするファイルのS3オブジェクトキー（例: "uploads/coco.png"）
//
// return:
//   - error: error情報
func (cs3 *s3Provider) GetObject(c echo.Context, sok string) error {
	logger := log.GetLogger(c).Sugar()

	logger.Debugf("Bucket: %v, Key: %v", configs.FetchConfigStr("aws.s3.bucket.name"), sok)

	// 取得したいs3への情報
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(configs.FetchConfigStr("aws.s3.bucket.name")),
		Key:    aws.String(sok),
	}

	logger.Infof("GetObject input: %v", *getObjectInput)

	// オプションを取得
	optFns := getS3Options()

	// s3の対象のオブジェクト取得
	getObjectOutput, err := cs3.svc.GetObject(
		context.Background(),
		getObjectInput,
		optFns...,
	)

	logger.Infof("GetObject output: %+v", getObjectOutput)

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"画像、ファイルの取得に失敗しました。",
			wrErrors.NewCmsServerErrorEType(),
		)

		logger.Errorf("S3 getObject failure: %v", wrErr)
		return wrErr
	}

	// logger.Info("Success s3 object delete!!!")

	return nil
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

// checkS3ObjectNotExists: S3にオブジェクトがないことを確認する関数(削除済みかの確認)
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - string: アップロードするファイルのS3オブジェクトキー（例: "uploads/coco.png"）
//
// return:
//   - error: error情報
func (cs3 *s3Provider) checkS3ObjectNotExists(c echo.Context, sok string) error {
	logger := log.GetLogger(c).Sugar()

	logger.Debugf("Bucket: %v, Key: %v", configs.FetchConfigStr("aws.s3.bucket.name"), sok)

	// 確認したい S3 オブジェクトの情報
	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(configs.FetchConfigStr("aws.s3.bucket.name")),
		Key:    aws.String(sok),
	}

	logger.Infof("HeadObject input: %+v", headObjectInput)

	// オプションを取得
	optFns := getS3Options()

	// S3 のオブジェクト存在確認
	headObjectOutput, err := cs3.svc.HeadObject(
		context.Background(),
		headObjectInput,
		optFns...,
	)

	logger.Infof("GetObject output: %+v", headObjectOutput)

	if err != nil {
		var apiErr smithy.APIError
		// 対象のオブジェクトがないので正常に削除されている
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
			logger.Infof("S3 object does not exist: %v. It has already been deleted or does not require deletion.", sok)
		} else {
			// 通常のエラー
			wrErr := wrErrors.NewWRError(
				err,
				"画像、ファイルの取得に失敗しました。",
				wrErrors.NewCmsServerErrorEType(),
			)

			logger.Errorf("S3 HeadObject failure: %v", wrErr)
			return wrErr
		}
	} else {
		logger.Warnf("S3 object exists: %v", sok)
	}

	return nil
}
