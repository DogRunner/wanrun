package handler

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/cms/adapters/aws"
	"github.com/wanrun-develop/wanrun/internal/cms/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/cms/core/dto"
	model "github.com/wanrun-develop/wanrun/internal/models"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
	wrUtil "github.com/wanrun-develop/wanrun/pkg/util"
)

const (
	S3_ROOT_FOLDER    = "cms"
	S3_SERVICE_FOLDER = "wanrun"
)

type ICmsHandler interface {
	HandleFileUpload(c echo.Context, fuq dto.FileUploadReq) (dto.FileUploadRes, error)
}

type cmsHandler struct {
	cs3 aws.IS3Provider
	cr  repository.ICmsRepository
}

func NewCmsHandler(cs3 aws.IS3Provider, cr repository.ICmsRepository) ICmsHandler {
	return &cmsHandler{cs3, cr}
}

// HandleFileUpload: S3へアップロードとDB登録
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - dto.FileUploadReq: フロントからのリクエス情報
//
// return:
//   - error: error情報
func (ch *cmsHandler) HandleFileUpload(c echo.Context, fuq dto.FileUploadReq) (dto.FileUploadRes, error) {
	// fileIDの生成
	fileID, wrErr := generateFileID(c)

	if wrErr != nil {
		return dto.FileUploadRes{}, wrErr
	}

	// s3オブジェクトキーの生成
	s3ObjectKey := generateS3ObjectKey(fileID, fuq)

	// s3へのアップロード
	if wrErr := ch.cs3.PutObject(c, s3ObjectKey, fuq.Src); wrErr != nil {
		return dto.FileUploadRes{}, wrErr
	}

	// fileのサイズ取得
	fileSize, wrErr := getFileSize(c, fuq.Src)
	if wrErr != nil {
		return dto.FileUploadRes{}, wrErr
	}

	s3FI := model.S3FileInfo{
		FileID:      wrUtil.NewSqlNullString(fileID),
		FileSize:    wrUtil.NewSqlNullInt64(fileSize),
		S3ObjectKey: wrUtil.NewSqlNullString(s3ObjectKey),
		DogownerID:  wrUtil.NewSqlNullInt64(fuq.DogownerID),
	}

	// S3FileInfoの登録
	if wrErr := ch.cr.CreateS3FileInfo(c, s3FI); wrErr != nil {
		return dto.FileUploadRes{}, wrErr
	}

	fuRes := dto.FileUploadRes{
		FileID: s3FI.FileID.String,
	}

	return fuRes, nil
}

// generateS3ObjectKey: S3ObjectKeyの生成
//
// args:
//   - string: 生成したfileID
//   - dto.FileUploadReq: フロントから来たfileUpload情報
//
// return:
//   - string: s3ObjectKey
func generateS3ObjectKey(fileID string, fuq dto.FileUploadReq) string {
	return fmt.Sprintf("%s/%s/%s/%s.%s",
		S3_ROOT_FOLDER,
		S3_SERVICE_FOLDER,
		fileID,
		fuq.FileName,
		fuq.Extension,
	)
}

// generateFileID: FileIDの生成。引数の数だけランダムの文字列を生成
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//
// return:
//   - string: fileID
//   - error: error情報
func generateFileID(c echo.Context) (string, error) {
	logger := log.GetLogger(c).Sugar()

	// カスタムエラー処理
	handleError := func(err error) error {
		wrErr := wrErrors.NewWRError(
			err,
			"FileID生成に失敗しました",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Error(wrErr)
		return wrErr
	}

	// UUIDを生成
	return util.UUIDGenerator(handleError)
}

// getFileSize: fileサイズの取得
//
// args:
//   - echo.Context: Echoのコンテキスト。リクエストやレスポンスにアクセスするために使用
//   - multipart.File: フロントで取得したファイルデータ
//
// return:
//   - int64: fileサイズ
//   - error: error情報
func getFileSize(c echo.Context, file multipart.File) (int64, error) {
	logger := log.GetLogger(c).Sugar()
	// 現在の位置を保存
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"Fileの位置サイズの取得に失敗しました",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Error(wrErr)
		return 0, wrErr
	}

	// ファイルのサイズを取得
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"Fileの取得に失敗しました",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Error(wrErr)
		return 0, wrErr
	}

	// 元の位置に戻す
	_, err = file.Seek(currentPos, io.SeekStart)
	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"Fileの位置を戻すのに失敗しました",
			wrErrors.NewCmsServerErrorEType(),
		)
		logger.Error(wrErr)
		return 0, wrErr
	}

	return size, nil
}
