package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockPostgres(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}
	return db, mock
}

func TestSaveOrderRollsBackWhenTradePointLookupFails(t *testing.T) {
	db, mock := newMockPostgres(t)
	repo := NewOrderPostgres(db)
	dbErr := errors.New("database unavailable")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "trade_points" WHERE "trade_points"."id" = $1 ORDER BY "trade_points"."id" LIMIT $2`)).
		WithArgs(10, 1).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.SaveOrder(9, 6, 2026, 10, 1, false, false, []models.CreateOrder{{Id: 5, Count: 2}})
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected database error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("transaction expectations: %v", err)
	}
}

func TestSavePaymentRollsBackWhenPaymentExceedsOrderSum(t *testing.T) {
	db, mock := newMockPostgres(t)
	repo := NewPaymentPostgres(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1`)).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sale_sum"}).AddRow(7, 100.0))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE order_id = $1`)).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "currency"}).AddRow(1, 7, 90.0))
	mock.ExpectRollback()

	err := repo.SavePayment(3, models.CreatePayment{
		OrderId:       7,
		PaymentTypeId: 1,
		Currency:      20,
	})
	if err == nil {
		t.Fatal("expected overpayment error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("transaction expectations: %v", err)
	}
}
