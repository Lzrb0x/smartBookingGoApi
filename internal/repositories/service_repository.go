package repositories

import (
	"context"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type ServiceRepository struct {
	db *database.DB
}

func NewServiceRepository(db *database.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) FindAll(ctx context.Context) ([]models.Service, error) {
	var services []models.Service
	err := r.db.SQL.SelectContext(ctx, &services, `SELECT * FROM services`)
	return services, err
}

func (r *ServiceRepository) FindByID(ctx context.Context, id int64) (*models.Service, error) {
	var service models.Service
	err := r.db.SQL.GetContext(ctx, &service, `SELECT * FROM services WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Create(ctx context.Context, service *models.Service) error {
	query := `
		INSERT INTO services (name, description)
		VALUES (:name, :description)
		RETURNING id`
	rows, err := r.db.SQL.NamedQueryContext(ctx, query, service)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&service.ID)
	}
	return nil
}

func (r *ServiceRepository) Update(ctx context.Context, service *models.Service) error {
	query := `
		UPDATE services
		SET name = :name, description = :description
		WHERE id = :id`
	_, err := r.db.SQL.NamedExecContext(ctx, query, service)
	return err
}

func (r *ServiceRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.SQL.ExecContext(ctx, `DELETE FROM services WHERE id = $1`, id)
	return err
}
