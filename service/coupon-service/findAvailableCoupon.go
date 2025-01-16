package coupon_service_detail

import (
	"database/sql"
	"errors"
)

func FindAvailableCoupon(tx *sql.Tx) (string, error) {
	var couponCode string
	err := tx.QueryRow("SELECT COUPON_CODE FROM COUPONS WHERE STATUS = 'AVAILABLE' LIMIT 1").Scan(&couponCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("쿠폰이 모두 소진되었습니다")
		}
		return "", err
	}
	return couponCode, nil
}
