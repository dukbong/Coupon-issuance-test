package coupon_service_detail

import "database/sql"

func UpdateCouponStatus(tx *sql.Tx, couponCode string, userId string) error {
	_, err := tx.Exec(
		"UPDATE coupons SET status = 'ISSUED', issued_to = ? WHERE coupon_code = ?",
		userId, couponCode,
	)
	return err
}
