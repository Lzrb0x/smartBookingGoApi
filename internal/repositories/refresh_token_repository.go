package repositories

import (
	"context"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type RefreshTokenRepository struct {
	db *database.DB
}

func NewRefreshTokenRepository(db *database.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_id, token_hash, expires_on)
		VALUES (:user_id, :token_id, :token_hash, :expires_on)
		RETURNING id, created_on`

	rows, err := r.db.SQL.NamedQueryContext(ctx, query, token)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&token.ID, &token.CreatedOn)
	}

	return nil
}

func (r *RefreshTokenRepository) FindActiveByTokenID(ctx context.Context, tokenID string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.SQL.GetContext(
		ctx,
		&token,
		`SELECT * FROM refresh_tokens WHERE token_id = $1 AND revoked_on IS NULL`,
		tokenID,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) RevokeByTokenID(ctx context.Context, tokenID string) error {
	_, err := r.db.SQL.ExecContext(
		ctx,
		`UPDATE refresh_tokens SET revoked_on = NOW() WHERE token_id = $1 AND revoked_on IS NULL`,
		tokenID,
	)
	return err
}
