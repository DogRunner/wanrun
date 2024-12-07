package dto

import "mime/multipart"

type FileUploadReq struct {
	FileName   string         // ファイル名
	Extension  string         // ファイルの拡張子 (例: ".png", ".txt")
	Src        multipart.File // ファイルの内容
	DogOwnerID int64
}
