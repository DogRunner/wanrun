package errors

import "fmt"

const (
	OTHER       int = 0
	AUTH        int = 1
	DOG         int = 2
	DOG_OWNER   int = 3
	DOGRUN      int = 4
	CMS         int = 5
	ORG         int = 6
	DOGRUNMG    int = 7
	INTERACTION int = 8
)

const (
	CLIENT int = 1
	SERVER int = 2
)

type eType struct {
	service   int
	errorType int
}

func (t eType) String() string {
	// カスタムフォーマットで文字列を返す
	return fmt.Sprintf("%d-%d", t.service, t.errorType)
}

/*
その他の予期せぬエラー
*/
func NewUnexpectedErrorEType() eType {
	return eType{OTHER, SERVER}
}

/*
認証機能のクライアントエラー
*/
func NewAuthClientErrorEType() eType {
	return eType{AUTH, CLIENT}
}

/*
認証機能のサーバーエラー
*/
func NewAuthServerErrorEType() eType {
	return eType{AUTH, SERVER}
}

/*
ドッグ機能のクライアントエラー
*/
func NewDogClientErrorEType() eType {
	return eType{DOG, CLIENT}
}

/*
ドッグ機能のサーバーエラー
*/
func NewDogServerErrorEType() eType {
	return eType{DOG, SERVER}
}

/*
ドッグオーナー機能のクライアントエラー
*/
func NewDogOwnerClientErrorEType() eType {
	return eType{DOG_OWNER, CLIENT}
}

/*
ドッグオーナー機能のサーバーエラー
*/
func NewDogOwnerServerErrorEType() eType {
	return eType{DOG_OWNER, SERVER}
}

/*
ドッグラン機能のクライアントエラー
*/
func NewDogrunClientErrorEType() eType {
	return eType{DOGRUN, CLIENT}
}

/*
ドッグラン機能のサーバーエラー
*/
func NewDogrunServerErrorEType() eType {
	return eType{DOGRUN, SERVER}
}

/*
インタラクション機能のクライアントエラー
*/
func NewInteractionClientErrorEType() eType {
	return eType{INTERACTION, CLIENT}
}

/*
インタラクション機能のサーバーエラー
*/
func NewInteractionServerErrorEType() eType {
	return eType{INTERACTION, SERVER}
}

/*
cms機能のクライアントエラー
*/
func NewCmsClientErrorEType() eType {
	return eType{CMS, CLIENT}
}

/*
cms機能のサーバーエラー
*/
func NewCmsServerErrorEType() eType {
	return eType{CMS, SERVER}
}

/*
organization機能のクライアントエラー
*/
func NewOrgClientErrorEType() eType {
	return eType{ORG, CLIENT}
}

/*
organization機能のサーバーエラー
*/
func NewOrgServerErrorEType() eType {
	return eType{ORG, SERVER}
}

/*
dogrunmg機能のクライアントエラー
*/
func NewDogrunmgClientErrorEType() eType {
	return eType{DOGRUNMG, CLIENT}
}

/*
dogrunmg機能のサーバーエラー
*/
func NewDogrunmgServerErrorEType() eType {
	return eType{DOGRUNMG, SERVER}
}
