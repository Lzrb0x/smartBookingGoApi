package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

func newMockDB(t *testing.T) (*database.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	cleanup := func() {
		_ = mockDB.Close()
	}

	return &database.DB{SQL: sqlxDB}, mock, cleanup
}

func TestBarbershopServiceRepository_Update_NoRows(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewBarbershopServiceRepository(db)

	service := &models.BarbershopService{
		ID:                  5,
		BarbershopID:        7,
		Price:               10.0,
		Duration:            30,
		DescriptionOverride: "",
		ServiceID:           2,
	}

	mock.ExpectExec(`UPDATE barbershop_services`).
		WithArgs(10.0, 30, "", int64(2), int64(5), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(context.Background(), service)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestBarbershopServiceRepository_Update_OK(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewBarbershopServiceRepository(db)

	service := &models.BarbershopService{
		ID:                  5,
		BarbershopID:        7,
		Price:               10.0,
		Duration:            30,
		DescriptionOverride: "",
		ServiceID:           2,
	}

	mock.ExpectExec(`UPDATE barbershop_services`).
		WithArgs(10.0, 30, "", int64(2), int64(5), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Update(context.Background(), service); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestBarbershopServiceRepository_Delete_NoRows(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewBarbershopServiceRepository(db)

	mock.ExpectExec(`DELETE FROM barbershop_services`).
		WithArgs(int64(5), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(context.Background(), 7, 5)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestBarbershopServiceRepository_Delete_OK(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewBarbershopServiceRepository(db)

	mock.ExpectExec(`DELETE FROM barbershop_services`).
		WithArgs(int64(5), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Delete(context.Background(), 7, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
