package repositories

import (
	"context"
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type EmployeeWorkingHourOverrideRepository struct {
	db *database.DB
}

func NewEmployeeWorkingHourOverrideRepository(db *database.DB) *EmployeeWorkingHourOverrideRepository {
	return &EmployeeWorkingHourOverrideRepository{db: db}
}

func (r *EmployeeWorkingHourOverrideRepository) Create(ctx context.Context, override *models.EmployeeWorkingHourOverride) error {
	query := `INSERT INTO employee_working_hour_overrides (employee_id, date, start_time, end_time, is_day_off) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query, override.EmployeeID, override.Date, override.StartTime, override.EndTime, override.IsDayOff).Scan(&override.ID)
}

func (r *EmployeeWorkingHourOverrideRepository) FindByID(ctx context.Context, id int64) (*models.EmployeeWorkingHourOverride, error) {
	var override models.EmployeeWorkingHourOverride
	err := r.db.SQL.GetContext(ctx, &override,
		`SELECT id, employee_id, date, start_time, end_time, is_day_off FROM employee_working_hour_overrides WHERE id = $1`, id)
	return &override, err
}

func (r *EmployeeWorkingHourOverrideRepository) FindByEmployee(ctx context.Context, employeeID int64) ([]models.EmployeeWorkingHourOverride, error) {
	var overrides []models.EmployeeWorkingHourOverride
	err := r.db.SQL.SelectContext(ctx, &overrides,
		`SELECT id, employee_id, date, start_time, end_time, is_day_off FROM employee_working_hour_overrides WHERE employee_id = $1 ORDER BY date DESC`, employeeID)
	return overrides, err
}

func (r *EmployeeWorkingHourOverrideRepository) FindByEmployeeAndBarbershop(ctx context.Context, employeeID, barbershopID int64) ([]models.EmployeeWorkingHourOverride, error) {
	var overrides []models.EmployeeWorkingHourOverride
	err := r.db.SQL.SelectContext(ctx, &overrides,
		`SELECT o.id, o.employee_id, o.date, o.start_time, o.end_time, o.is_day_off 
		 FROM employee_working_hour_overrides o
		 JOIN employees e ON o.employee_id = e.id
		 WHERE o.employee_id = $1 AND e.barbershop_id = $2 
		 ORDER BY o.date DESC`, employeeID, barbershopID)
	return overrides, err
}

func (r *EmployeeWorkingHourOverrideRepository) Update(ctx context.Context, override *models.EmployeeWorkingHourOverride) error {
	_, err := r.db.SQL.ExecContext(ctx,
		`UPDATE employee_working_hour_overrides SET date = $1, start_time = $2, end_time = $3, is_day_off = $4 WHERE id = $5`,
		override.Date, override.StartTime, override.EndTime, override.IsDayOff, override.ID)
	return err
}

func (r *EmployeeWorkingHourOverrideRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.SQL.ExecContext(ctx, `DELETE FROM employee_working_hour_overrides WHERE id = $1`, id)
	return err
}

func (r *EmployeeWorkingHourOverrideRepository) ExistsByEmployeeAndDate(ctx context.Context, employeeID int64, date time.Time) (bool, error) {
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM employee_working_hour_overrides WHERE employee_id = $1 AND date = $2)`,
		employeeID, date).Scan(&exists)
	return exists, err
}

func (r *EmployeeWorkingHourOverrideRepository) ExistsByEmployeeAndDateExcluding(ctx context.Context, employeeID int64, date time.Time, excludeID int64) (bool, error) {
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM employee_working_hour_overrides WHERE employee_id = $1 AND date = $2 AND id != $3)`,
		employeeID, date, excludeID).Scan(&exists)
	return exists, err
}
