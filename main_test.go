package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	CouponMainService "coupon-server/service"

	_ "github.com/glebarez/sqlite"
)

// 테스트를 위한 메모리 DB 설정
func setupTestDB() (*sql.DB, error) {
	// 메모리 DB로 SQLite 연결
	db, err := sql.Open("sqlite", ":memory:") // 메모리 DB 사용
	if err != nil {
		return nil, fmt.Errorf("DB 연결 실패: %w", err)
	}

	// DB 연결 상태 확인
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 상태 확인 실패: %w", err)
	}

	// 메모리 DB 테이블 생성
	createTableQuery := `
	CREATE TABLE COUPONS (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		COUPON_CODE TEXT NOT NULL UNIQUE,
		STATUS TEXT NOT NULL,
		ISSUED_TO TEXT,
		ISSUED_AT DATETIME
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("테이블 생성 실패: %w", err)
	}

	insertQuery := `
		INSERT INTO COUPONS (COUPON_CODE, STATUS)
		VALUES (?, ?)
	`

	for i := 0; i < 15; i++ {
		_, err := db.Exec(insertQuery, fmt.Sprintf("coupon%d", i), "AVAILABLE")
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func TestIssueCouponConcurrently(t *testing.T) {
	// 테스트를 위한 DB 설정
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("DB 설정 실패: %v", err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	numRequests := 10000
	startTime := time.Now()
	// 동시 요청 테스트
	for i := 0; i < numRequests; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			userId := fmt.Sprintf("user-%d", i)
			err := CouponMainService.IssueCoupon(db, userId)
			if err != nil {
				t.Errorf("user-%d 쿠폰 발급 실패: %v", i, err)
			}
		}(i)
	}

	// 모든 고루틴이 종료될 때까지 기다림
	wg.Wait()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("선착순 쿠폰 발급 처리 시간: %v\n", duration)
	// 발급된 쿠폰 상태 확인
	type QueryResult struct {
		CouponCode string
		IssuedTo   string
		Status     string
	}
	rows, err := db.Query("SELECT COUPON_CODE, ISSUED_TO, STATUS FROM COUPONS WHERE STATUS = 'ISSUED'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // rows를 다 사용한 후 반드시 Close 해줘야 합니다.

	// 결과 처리
	for rows.Next() {
		var coupon QueryResult
		if err := rows.Scan(&coupon.CouponCode, &coupon.IssuedTo, &coupon.Status); err != nil {
			log.Fatal(err)
		}
		// 각 쿠폰 정보를 출력
		fmt.Printf("쿠폰 코드: %s, 발급 대상: %s, 상태: %s\n", coupon.CouponCode, coupon.IssuedTo, coupon.Status)
	}

	// rows.Next() 후 오류 처리
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
