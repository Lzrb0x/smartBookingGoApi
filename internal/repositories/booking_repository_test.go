package repositories

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func professionalBookingRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id",
		"customer_id",
		"customer_name",
		"customer_phone",
		"employee_id",
		"employee_user_id",
		"employee_name",
		"employee_phone",
		"barbershop_id",
		"barbershop_name",
		"barbershop_service_id",
		"service_id",
		"service_name",
		"service_price",
		"service_duration",
		"date",
		"start_time",
		"end_time",
	}).AddRow(
		int64(1),
		int64(10),
		"Cliente",
		"85999999999",
		int64(20),
		int64(30),
		"Barbeiro",
		"85888888888",
		int64(40),
		"Barbearia",
		int64(50),
		int64(60),
		"Corte",
		35.5,
		30,
		time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC),
		time.Date(0000, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(0000, 1, 1, 9, 30, 0, 0, time.UTC),
	)
}

func TestBookingRepository_ListProfessionalDashboard_AllEmployees(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewBookingRepository(db)
	date := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`FROM bookings b`).
		WithArgs(int64(40), date).
		WillReturnRows(professionalBookingRows())

	bookings, err := repo.ListProfessionalDashboard(context.Background(), 40, date, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	if bookings[0].CustomerName != "Cliente" || bookings[0].ServicePrice != 35.5 {
		t.Fatalf("unexpected booking payload: %+v", bookings[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestBookingRepository_ListProfessionalDashboard_ByEmployee(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewBookingRepository(db)
	date := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	employeeID := int64(20)

	mock.ExpectQuery(`AND b.employee_id = \$3`).
		WithArgs(int64(40), date, employeeID).
		WillReturnRows(professionalBookingRows())

	bookings, err := repo.ListProfessionalDashboard(context.Background(), 40, date, &employeeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bookings) != 1 || bookings[0].EmployeeID != employeeID {
		t.Fatalf("expected employee-filtered booking, got %+v", bookings)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
