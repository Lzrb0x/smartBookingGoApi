package repositories

import (
	"context"
	"database/sql"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type EmployeeRepository struct {
	db *database.DB
}

func NewEmployeeRepository(db *database.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) FindByBarbershop(ctx context.Context, barbershopID int64) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.SQL.SelectContext(ctx, &employees,
		`SELECT * FROM employees WHERE barbershop_id = $1`, barbershopID)
	return employees, err
}

func (r *EmployeeRepository) FindByUser(ctx context.Context, userID int64) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.SQL.SelectContext(ctx, &employees,
		`SELECT * FROM employees WHERE user_id = $1 ORDER BY barbershop_id, id`, userID)
	return employees, err
}

func (r *EmployeeRepository) FindByUserAndBarbershop(ctx context.Context, userID, barbershopID int64) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.SQL.GetContext(ctx, &employee,
		`SELECT * FROM employees WHERE user_id = $1 AND barbershop_id = $2`, userID, barbershopID)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	existing, err := r.FindByUserAndBarbershop(ctx, employee.UserID, employee.BarberShopID)
	if err == nil {
		employee.ID = existing.ID
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}

	query := `INSERT INTO employees (user_id, barbershop_id) VALUES ($1, $2) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query, employee.UserID, employee.BarberShopID).Scan(&employee.ID)
}

func (r *EmployeeRepository) Delete(ctx context.Context, barbershopID, employeeID int64) error {
	_, err := r.db.SQL.ExecContext(ctx,
		`DELETE FROM employees WHERE id = $1 AND barbershop_id = $2`, employeeID, barbershopID)
	return err
}

func (r *EmployeeRepository) ExistsInBarbershop(ctx context.Context, employeeID, barbershopID int64) (bool, error) {
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM employees WHERE id = $1 AND barbershop_id = $2)`, employeeID, barbershopID).Scan(&exists)
	return exists, err
}
