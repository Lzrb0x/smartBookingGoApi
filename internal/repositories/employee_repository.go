package repositories

import (
	"context"

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

func (r *EmployeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	query := `INSERT INTO employees (user_id, barbershop_id) VALUES ($1, $2) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query, employee.UserID, employee.BarberShopID).Scan(&employee.ID)
}

func (r *EmployeeRepository) Delete(ctx context.Context, barbershopID, employeeID int64) error {
	_, err := r.db.SQL.ExecContext(ctx,
		`DELETE FROM employees WHERE id = $1 AND barbershop_id = $2`, employeeID, barbershopID)
	return err
}
