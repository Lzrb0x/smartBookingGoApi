package repositories

import (
	"context"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type ServiceEmployeeRepository struct {
	db *database.DB
}

func NewServiceEmployeeRepository(db *database.DB) *ServiceEmployeeRepository {
	return &ServiceEmployeeRepository{db: db}
}

func (r *ServiceEmployeeRepository) AssignServiceToEmployee(ctx context.Context, serviceEmployee *models.ServiceEmployee) error {
	query := `INSERT INTO service_employees (employee_id, barbershop_service_id) VALUES ($1, $2) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query, serviceEmployee.EmployeeID, serviceEmployee.BarbershopServiceID).Scan(&serviceEmployee.ID)
}

func (r *ServiceEmployeeRepository) UnassignServiceFromEmployee(ctx context.Context, employeeID, barbershopServiceID int64) error {
	_, err := r.db.SQL.ExecContext(ctx,
		`DELETE FROM service_employees WHERE employee_id = $1 AND barbershop_service_id = $2`, employeeID, barbershopServiceID)
	return err
}

func (r *ServiceEmployeeRepository) GetServicesByEmployee(ctx context.Context, employeeID int64) ([]models.ServiceEmployee, error) {
	var serviceEmployees []models.ServiceEmployee
	err := r.db.SQL.SelectContext(ctx, &serviceEmployees,
		`SELECT * FROM service_employees WHERE employee_id = $1 ORDER BY id`, employeeID)
	return serviceEmployees, err
}

func (r *ServiceEmployeeRepository) GetEmployeesByService(ctx context.Context, barbershopServiceID int64) ([]models.ServiceEmployee, error) {
	var serviceEmployees []models.ServiceEmployee
	err := r.db.SQL.SelectContext(ctx, &serviceEmployees,
		`SELECT * FROM service_employees WHERE barbershop_service_id = $1 ORDER BY id`, barbershopServiceID)
	return serviceEmployees, err
}

func (r *ServiceEmployeeRepository) IsEmployeeAssignedToService(ctx context.Context, employeeID, barbershopServiceID int64) (bool, error) {
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM service_employees WHERE employee_id = $1 AND barbershop_service_id = $2)`, employeeID, barbershopServiceID).Scan(&exists)
	return exists, err
}

func (r *ServiceEmployeeRepository) BarbershopServiceExists(ctx context.Context, barbershopServiceID, barbershopID int64) (bool, error) {
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM barbershop_services WHERE id = $1 AND barbershop_id = $2)`, barbershopServiceID, barbershopID).Scan(&exists)
	return exists, err
}

func (r *ServiceEmployeeRepository) EmployeeExistsInBarbershop(ctx context.Context, employeeID, barbershopID int64) (bool, error) {
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM employees WHERE id = $1 AND barbershop_id = $2)`, employeeID, barbershopID).Scan(&exists)
	return exists, err
}
