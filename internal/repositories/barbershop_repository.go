package repositories

import (
	"context"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type BarbershopRepository struct {
	db *database.DB
}

func NewBarbershopRepository(db *database.DB) *BarbershopRepository {
	return &BarbershopRepository{db: db}
}

func (r *BarbershopRepository) FindAll(ctx context.Context) ([]models.Barbershop, error) {
	var barbershops []models.Barbershop
	err := r.db.SQL.SelectContext(ctx, &barbershops, `SELECT * FROM barbershops`)
	return barbershops, err
}

func (r *BarbershopRepository) FindByID(ctx context.Context, id int64) (*models.Barbershop, error) {
	var barbershop models.Barbershop
	err := r.db.SQL.GetContext(ctx, &barbershop, `SELECT * FROM barbershops WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &barbershop, nil
}

func (r *BarbershopRepository) Create(ctx context.Context, barbershop *models.Barbershop) error {
	query := `
		INSERT INTO barbershops (barbershop_name, address, phone, owner_id)
		VALUES (:barbershop_name, :address, :phone, :owner_id)
		RETURNING id`
	rows, err := r.db.SQL.NamedQueryContext(ctx, query, barbershop)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&barbershop.ID)
	}
	return nil
}

func (r *BarbershopRepository) Update(ctx context.Context, barbershop *models.Barbershop) error {
	query := `
		UPDATE barbershops
		SET barbershop_name = :barbershop_name, address = :address, phone = :phone
		WHERE id = :id`
	_, err := r.db.SQL.NamedExecContext(ctx, query, barbershop)
	return err
}

func (r *BarbershopRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.SQL.ExecContext(ctx, `DELETE FROM barbershops WHERE id = $1`, id)
	return err
}
