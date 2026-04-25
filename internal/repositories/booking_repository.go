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

func (r *BookingRepository) ListRecentByCustomer(ctx context.Context, customerID int64, limit int) ([]models.BookingDashboardItem, error) {
	var bookings []models.BookingDashboardItem
	err := r.db.SQL.SelectContext(ctx, &bookings, bookingDashboardSelect(`
		WHERE b.customer_id = $1
		ORDER BY b.date DESC, b.start_time DESC
		LIMIT $2`), customerID, limit)
	return bookings, err
}

func (r *BookingRepository) FindCurrentByCustomer(ctx context.Context, customerID int64) (*models.BookingDashboardItem, error) {
	var booking models.BookingDashboardItem
	err := r.db.SQL.GetContext(ctx, &booking, bookingDashboardSelect(`
		WHERE b.customer_id = $1
		  AND (b.date > CURRENT_DATE OR (b.date = CURRENT_DATE AND b.end_time >= CURRENT_TIME))
		ORDER BY b.date ASC, b.start_time ASC
		LIMIT 1`), customerID)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) ListRecentBarbershopsByCustomer(ctx context.Context, customerID int64, limit int) ([]models.RecentBarbershop, error) {
	var barbershops []models.RecentBarbershop
	err := r.db.SQL.SelectContext(ctx, &barbershops, `
		SELECT id, barbershop_name, address, phone, last_booking_date, last_start_time
		FROM (
			SELECT DISTINCT ON (ba.id)
				ba.id,
				ba.barbershop_name,
				COALESCE(ba.address, '') AS address,
				COALESCE(ba.phone, '') AS phone,
				b.date AS last_booking_date,
				b.start_time AS last_start_time
			FROM bookings b
			INNER JOIN barbershops ba ON ba.id = b.barbershop_id
			WHERE b.customer_id = $1
			ORDER BY ba.id, b.date DESC, b.start_time DESC
		) recent_barbershops
		ORDER BY last_booking_date DESC, last_start_time DESC
		LIMIT $2`, customerID, limit)
	return barbershops, err
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

func bookingDashboardSelect(suffix string) string {
	return `
		SELECT
			b.id,
			b.customer_id,
			b.employee_id,
			COALESCE(employee_user.name, '') AS employee_name,
			b.barbershop_id,
			ba.barbershop_name,
			COALESCE(ba.address, '') AS barbershop_address,
			COALESCE(ba.phone, '') AS barbershop_phone,
			b.barbershop_service_id,
			s.id AS service_id,
			s.name AS service_name,
			COALESCE(s.description, '') AS service_description,
			bs.price::float8 AS service_price,
			bs.duration AS service_duration,
			COALESCE(bs.description_override, '') AS service_description_override,
			b.date,
			b.start_time,
			b.end_time
		FROM bookings b
		INNER JOIN employees e ON e.id = b.employee_id
		INNER JOIN users employee_user ON employee_user.id = e.user_id
		INNER JOIN barbershops ba ON ba.id = b.barbershop_id
		INNER JOIN barbershop_services bs ON bs.id = b.barbershop_service_id
		INNER JOIN services s ON s.id = bs.service_id
	` + suffix
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
