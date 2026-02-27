package repositories

import (
	"context"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type OwnerRepository struct {
	db *database.DB
}

func NewOwnerRepository(db *database.DB) *OwnerRepository {
	return &OwnerRepository{db: db}
}

func (r *OwnerRepository) Create(ctx context.Context, owner *models.Owner) error {
	query := `INSERT INTO owners (user_id) VALUES ($1) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query, owner.UserID).Scan(&owner.ID)
}
