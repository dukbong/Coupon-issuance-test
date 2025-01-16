package coupon_service

import "database/sql"

// 트랜잭션 커밋 또는 롤백 함수
func CommitOrRollback(tx *sql.Tx, err *error) {
	if *err != nil {
		tx.Rollback()
	} else {
		*err = tx.Commit()
	}
}
