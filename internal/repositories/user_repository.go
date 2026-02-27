package repositories

import (
	"context"

	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.SQL.SelectContext(ctx, &users, `SELECT * FROM users WHERE active = true`)
	return users, err
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.db.SQL.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1 AND active = true`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (user_identifier, name, email, password, phone, is_complete)
		VALUES (:user_identifier, :name, :email, :password, :phone, :is_complete)
		RETURNING id, active, created_on`
	rows, err := r.db.SQL.NamedQueryContext(ctx, query, user)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&user.ID, &user.Active, &user.CreatedOn)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = :name, email = :email, phone = :phone, is_complete = :is_complete
		WHERE id = :id AND active = true`
	_, err := r.db.SQL.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.SQL.ExecContext(ctx, `UPDATE users SET active = false WHERE id = $1`, id)
	return err
}
