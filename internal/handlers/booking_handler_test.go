package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func newBookingHandlerMock(t *testing.T) (*BookingHandler, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db := &database.DB{SQL: sqlx.NewDb(mockDB, "sqlmock")}
	handler := NewBookingHandler(
		repositories.NewBookingRepository(db),
		repositories.NewEmployeeRepository(db),
		repositories.NewOwnerRepository(db),
		repositories.NewBarbershopRepository(db),
		repositories.NewBarbershopServiceRepository(db),
		repositories.NewServiceEmployeeRepository(db),
		repositories.NewEmployeeWorkingHourRepository(db),
		repositories.NewEmployeeWorkingHourOverrideRepository(db),
	)

	return handler, mock, func() {
		_ = mockDB.Close()
	}
}

func newProfessionalContext(query string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest(http.MethodGet, "/?"+query, nil)
	context.Request = request
	return context, recorder
}

func expectBarbershop(mock sqlmock.Sqlmock, barbershopID, ownerID int64) {
	mock.ExpectQuery(`SELECT \* FROM barbershops WHERE id = \$1`).
		WithArgs(barbershopID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "barbershop_name", "address", "phone", "owner_id",
		}).AddRow(barbershopID, "Barbearia", "Rua A", "85999999999", ownerID))
}

func expectOwner(mock sqlmock.Sqlmock, userID, ownerID int64) {
	mock.ExpectQuery(`SELECT id, user_id FROM owners WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(ownerID, userID))
}

func expectNoOwner(mock sqlmock.Sqlmock, userID int64) {
	mock.ExpectQuery(`SELECT id, user_id FROM owners WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)
}

func expectEmployee(mock sqlmock.Sqlmock, userID, barbershopID, employeeID int64) {
	mock.ExpectQuery(`SELECT \* FROM employees WHERE user_id = \$1 AND barbershop_id = \$2`).
		WithArgs(userID, barbershopID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "barbershop_id"}).
			AddRow(employeeID, userID, barbershopID))
}

func expectNoEmployee(mock sqlmock.Sqlmock, userID, barbershopID int64) {
	mock.ExpectQuery(`SELECT \* FROM employees WHERE user_id = \$1 AND barbershop_id = \$2`).
		WithArgs(userID, barbershopID).
		WillReturnError(sql.ErrNoRows)
}

func TestAuthorizedProfessionalEmployeeFilter_OwnerCanSeeAll(t *testing.T) {
	handler, mock, cleanup := newBookingHandlerMock(t)
	defer cleanup()

	context, recorder := newProfessionalContext("employee_id=all")
	expectBarbershop(mock, 7, 3)
	expectOwner(mock, 10, 3)
	expectNoEmployee(mock, 10, 7)

	filter, ok := handler.authorizedProfessionalEmployeeFilter(context, 10, 7)
	if !ok || filter != nil {
		t.Fatalf("expected owner all-employees access, got ok=%v filter=%v status=%d", ok, filter, recorder.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestAuthorizedProfessionalEmployeeFilter_EmployeeFallsBackToOwnID(t *testing.T) {
	handler, mock, cleanup := newBookingHandlerMock(t)
	defer cleanup()

	context, recorder := newProfessionalContext("employee_id=all")
	expectBarbershop(mock, 7, 3)
	expectNoOwner(mock, 10)
	expectEmployee(mock, 10, 7, 22)

	filter, ok := handler.authorizedProfessionalEmployeeFilter(context, 10, 7)
	if !ok || filter == nil || *filter != 22 {
		t.Fatalf("expected employee own filter, got ok=%v filter=%v status=%d", ok, filter, recorder.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestAuthorizedProfessionalEmployeeFilter_RejectsUnlinkedUser(t *testing.T) {
	handler, mock, cleanup := newBookingHandlerMock(t)
	defer cleanup()

	context, recorder := newProfessionalContext("employee_id=all")
	expectBarbershop(mock, 7, 3)
	expectNoOwner(mock, 10)
	expectNoEmployee(mock, 10, 7)

	filter, ok := handler.authorizedProfessionalEmployeeFilter(context, 10, 7)
	if ok || filter != nil {
		t.Fatalf("expected forbidden access, got ok=%v filter=%v", ok, filter)
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", recorder.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
