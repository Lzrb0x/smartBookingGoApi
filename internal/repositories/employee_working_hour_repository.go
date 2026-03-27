package repositories

import (
	"context"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type EmployeeWorkingHourRepository struct {
	db *database.DB
}

func NewEmployeeWorkingHourRepository(db *database.DB) *EmployeeWorkingHourRepository {
	return &EmployeeWorkingHourRepository{db: db}
}

func (r *EmployeeWorkingHourRepository) Create(ctx context.Context, wh *models.EmployeeWorkingHour) error {
	query := `INSERT INTO employee_working_hours (employee_id, day_of_week, start_time, end_time, is_day_off) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query, wh.EmployeeID, wh.DayOfWeek, wh.StartTime, wh.EndTime, wh.IsDayOff).Scan(&wh.ID)
}

func (r *EmployeeWorkingHourRepository) FindByID(ctx context.Context, id int64) (*models.EmployeeWorkingHour, error) {
	var wh models.EmployeeWorkingHour
	err := r.db.SQL.GetContext(ctx, &wh,
		`SELECT id, employee_id, day_of_week, start_time, end_time, is_day_off FROM employee_working_hours WHERE id = $1`, id)
	return &wh, err
}

func (r *EmployeeWorkingHourRepository) FindByEmployee(ctx context.Context, employeeID int64) ([]models.EmployeeWorkingHour, error) {
	var whs []models.EmployeeWorkingHour
	err := r.db.SQL.SelectContext(ctx, &whs,
		`SELECT id, employee_id, day_of_week, start_time, end_time, is_day_off FROM employee_working_hours WHERE employee_id = $1 ORDER BY day_of_week`, employeeID)
	return whs, err
}

func (r *EmployeeWorkingHourRepository) FindByEmployeeAndBarbershop(ctx context.Context, employeeID, barbershopID int64) ([]models.EmployeeWorkingHour, error) {
	var whs []models.EmployeeWorkingHour
	err := r.db.SQL.SelectContext(ctx, &whs,
		`SELECT wh.id, wh.employee_id, wh.day_of_week, wh.start_time, wh.end_time, wh.is_day_off 
		 FROM employee_working_hours wh
		 JOIN employees e ON wh.employee_id = e.id
		 WHERE wh.employee_id = $1 AND e.barbershop_id = $2 
		 ORDER BY wh.day_of_week`, employeeID, barbershopID)
	return whs, err
}

func (r *EmployeeWorkingHourRepository) Update(ctx context.Context, wh *models.EmployeeWorkingHour) error {
	_, err := r.db.SQL.ExecContext(ctx,
		`UPDATE employee_working_hours SET day_of_week = $1, start_time = $2, end_time = $3, is_day_off = $4 WHERE id = $5`,
		wh.DayOfWeek, wh.StartTime, wh.EndTime, wh.IsDayOff, wh.ID)
	return err
}

func (r *EmployeeWorkingHourRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.SQL.ExecContext(ctx, `DELETE FROM employee_working_hours WHERE id = $1`, id)
	return err
}
