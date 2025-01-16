package coupon_service

import (
	"database/sql"
	"fmt"
)

// 트랜잭션 시작 함수
func BeginTransaction(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("트랜잭션 시작 실패: %w", err)
	}
	return tx, nil
}
