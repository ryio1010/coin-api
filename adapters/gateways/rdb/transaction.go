package rdb

import (
	"coin-api/domain/repository"
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

var txKey = struct{}{}

type TxRepository struct {
	DB *gorm.DB
}

func (tr *TxRepository) GetDBConn() *gorm.DB {
	return tr.DB
}

func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(&txKey).(*gorm.DB)
	return tx, ok
}

func NewTxRepository(DB *gorm.DB) repository.ITxRepository {
	return &TxRepository{
		DB: DB,
	}
}

func (tr *TxRepository) DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	// txを生成する
	conn := tr.GetDBConn()
	tx := conn.Begin(&sql.TxOptions{})

	ctx = context.WithValue(ctx, &txKey, tx)

	v, err := f(ctx)
	// エラーがあればロールバック
	if err != nil {
		_ = tx.Rollback()
		return v, fmt.Errorf("rollback: %w", err)
	}
	// エラーがなければコミット
	tx.Commit()
	return v, nil
}
