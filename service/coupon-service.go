package coupon_service

import (
	"database/sql"
	"fmt"
	"sync"

	CouponFunc "coupon-server/service/coupon-service"
)

var couponMutex sync.Mutex

func IssueCoupon(db *sql.DB, userId string) error {
	couponMutex.Lock()
	defer couponMutex.Unlock()

	tx, err := BeginTransaction(db)
	if err != nil {
		return err
	}

	defer CommitOrRollback(tx, &err)

	couponCode, err := CouponFunc.FindAvailableCoupon(tx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = CouponFunc.UpdateCouponStatus(tx, couponCode, userId)
	if err != nil {
		return fmt.Errorf("쿠폰 상태 업데이트 실패: %w", err)
	}

	return nil
}
