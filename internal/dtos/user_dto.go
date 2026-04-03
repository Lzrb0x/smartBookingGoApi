package dtos

import (
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name"     binding:"required"`
	Phone    string `json:"phone"    binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) ToModel() *models.User {
	var email *string
	if r.Email != "" {
		email = &r.Email
	}
	return &models.User{
		UserIdentifier: uuid.New().String(),
		Name:           r.Name,
		Phone:          r.Phone,
		Email:          email,
		Password:       r.Password,
	}
}
