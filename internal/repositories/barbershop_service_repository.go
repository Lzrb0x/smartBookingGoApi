package repositories

import (
	"context"
	"database/sql"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type BarbershopServiceRepository struct {
	db *database.DB
}

func NewBarbershopServiceRepository(db *database.DB) *BarbershopServiceRepository {
	return &BarbershopServiceRepository{db: db}
}

func (r *BarbershopServiceRepository) FindByBarbershop(ctx context.Context, barbershopID int64) ([]models.BarbershopService, error) {
	var services []models.BarbershopService
	err := r.db.SQL.SelectContext(ctx, &services,
		`SELECT * FROM barbershop_services WHERE barbershop_id = $1`, barbershopID)
	return services, err
}

func (r *BarbershopServiceRepository) FindByID(ctx context.Context, barbershopID, id int64) (*models.BarbershopService, error) {
	var service models.BarbershopService
	err := r.db.SQL.GetContext(ctx, &service,
		`SELECT * FROM barbershop_services WHERE id = $1 AND barbershop_id = $2`, id, barbershopID)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *BarbershopServiceRepository) Create(ctx context.Context, service *models.BarbershopService) error {
	query := `
		INSERT INTO barbershop_services (price, duration, description_override, barbershop_id, service_id)
		VALUES (:price, :duration, :description_override, :barbershop_id, :service_id)
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

func (r *BarbershopServiceRepository) Update(ctx context.Context, service *models.BarbershopService) error {
	query := `
		UPDATE barbershop_services
		SET price = :price, duration = :duration, description_override = :description_override, service_id = :service_id
		WHERE id = :id AND barbershop_id = :barbershop_id`
	result, err := r.db.SQL.NamedExecContext(ctx, query, service)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *BarbershopServiceRepository) Delete(ctx context.Context, barbershopID, id int64) error {
	result, err := r.db.SQL.ExecContext(ctx,
		`DELETE FROM barbershop_services WHERE id = $1 AND barbershop_id = $2`, id, barbershopID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
