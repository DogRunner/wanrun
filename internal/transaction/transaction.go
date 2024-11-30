package transaction

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

type TxKey struct{}

type ITransactionManager interface {
	Begin(ctx context.Context) (context.Context, *gorm.DB, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
	GetTx(ctx context.Context) *gorm.DB
	DoInTransaction(
		c echo.Context,
		ctx context.Context,
		f func(tx *gorm.DB) error,
		errorGenerator func(error) error, // エラー生成関数を受け取る
	) error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) ITransactionManager {
	return &transactionManager{db}
}

// トランザクションを開始し、contextに保存
func (tm *transactionManager) Begin(ctx context.Context) (context.Context, *gorm.DB, error) {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return ctx, nil, tx.Error
	}
	return context.WithValue(ctx, TxKey{}, tx), tx, nil
}

// トランザクションをコミット
func (tm *transactionManager) Commit(ctx context.Context) error {
	tx := ctx.Value(TxKey{}).(*gorm.DB)
	if tx == nil {
		return nil // トランザクションがない場合は何もしない
	}
	return tx.Commit().Error
}

// トランザクションをロールバック
func (tm *transactionManager) Rollback(ctx context.Context) {
	tx := ctx.Value(TxKey{}).(*gorm.DB)
	if tx != nil {
		tx.Rollback()
	}
}

// コンテキストからトランザクションを取得
func (tm *transactionManager) GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(TxKey{}).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

// ヘルパー関数: トランザクションを簡略化
func (tm *transactionManager) DoInTransaction(
	c echo.Context,
	ctx context.Context,
	f func(tx *gorm.DB) error,
	errorGenerator func(error) error, // エラー生成関数を受け取る
) error {
	logger := log.GetLogger(c).Sugar()

	// トランザクション開始
	tx := tm.db.Begin()
	if tx.Error != nil {
		return errorGenerator(tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("エラーが発生したため、ロールバックします。詳細情報: %v", r)
			tx.Rollback()
		}
	}()

	// コールバック関数の実行
	if err := f(tx); err != nil {
		tx.Rollback()
		return errorGenerator(tx.Error)
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		return errorGenerator(tx.Error)
	}

	return nil
}
