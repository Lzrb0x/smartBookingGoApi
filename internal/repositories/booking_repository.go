package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
	"github.com/jmoiron/sqlx"
)

type BookingRepository struct {
	db *database.DB
}

func NewBookingRepository(db *database.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.SQL.BeginTxx(ctx, nil)
}

func (r *BookingRepository) ListByEmployeeAndDate(ctx context.Context, employeeID int64, date time.Time) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.SQL.SelectContext(ctx, &bookings,
		`SELECT id, customer_id, employee_id, barbershop_id, barbershop_service_id, date, start_time, end_time
		 FROM bookings WHERE employee_id = $1 AND date = $2 ORDER BY start_time`, employeeID, date)
	return bookings, err
}

func (r *BookingRepository) ListByBarbershop(ctx context.Context, barbershopID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.SQL.SelectContext(ctx, &bookings,
		`SELECT id, customer_id, employee_id, barbershop_id, barbershop_service_id, date, start_time, end_time
		 FROM bookings WHERE barbershop_id = $1 ORDER BY date DESC, start_time ASC`, barbershopID)
	return bookings, err
}

func (r *BookingRepository) ListByEmployeeAndBarbershop(ctx context.Context, employeeID, barbershopID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.SQL.SelectContext(ctx, &bookings,
		`SELECT id, customer_id, employee_id, barbershop_id, barbershop_service_id, date, start_time, end_time
		 FROM bookings WHERE employee_id = $1 AND barbershop_id = $2 ORDER BY date DESC, start_time ASC`, employeeID, barbershopID)
	return bookings, err
}

func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	query := `INSERT INTO bookings (customer_id, employee_id, barbershop_id, barbershop_service_id, date, start_time, end_time)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query,
		booking.CustomerID,
		booking.EmployeeID,
		booking.BarbershopID,
		booking.BarbershopServiceID,
		booking.Date,
		booking.StartTime,
		booking.EndTime,
	).Scan(&booking.ID)
}

func (r *BookingRepository) LockEmployee(ctx context.Context, tx *sqlx.Tx, employeeID int64) error {
	var id int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM employees WHERE id = $1 FOR UPDATE`, employeeID).Scan(&id)
	if err == sql.ErrNoRows {
		return err
	}
	return err
}

func (r *BookingRepository) HasOverlapTx(ctx context.Context, tx *sqlx.Tx, employeeID int64, date time.Time, startTime, endTime time.Time) (bool, error) {
	var exists bool
	err := tx.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM bookings
			WHERE employee_id = $1 AND date = $2
			AND NOT (end_time <= $3 OR start_time >= $4)
		)`, employeeID, date, startTime, endTime).Scan(&exists)
	return exists, err
}

func (r *BookingRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, booking *models.Booking) error {
	query := `INSERT INTO bookings (customer_id, employee_id, barbershop_id, barbershop_service_id, date, start_time, end_time)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return tx.QueryRowContext(ctx, query,
		booking.CustomerID,
		booking.EmployeeID,
		booking.BarbershopID,
		booking.BarbershopServiceID,
		booking.Date,
		booking.StartTime,
		booking.EndTime,
	).Scan(&booking.ID)
}
