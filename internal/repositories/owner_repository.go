package repositories

import (
	"context"
	"database/sql"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type OwnerRepository struct {
	db *database.DB
}

func NewOwnerRepository(db *database.DB) *OwnerRepository {
	return &OwnerRepository{db: db}
}

func (r *OwnerRepository) FindByUserID(ctx context.Context, userID int64) (*models.Owner, error) {
	var owner models.Owner
	err := r.db.SQL.GetContext(ctx, &owner, `SELECT id, user_id FROM owners WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *OwnerRepository) Create(ctx context.Context, owner *models.Owner) error {
	existing, err := r.FindByUserID(ctx, owner.UserID)
	if err == nil {
		owner.ID = existing.ID
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}

	query := `INSERT INTO owners (user_id) VALUES ($1) RETURNING id`
	return r.db.SQL.QueryRowContext(ctx, query, owner.UserID).Scan(&owner.ID)
}
