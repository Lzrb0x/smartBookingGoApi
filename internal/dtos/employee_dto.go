package dtos

import "github.com/Lzrb0x/smartBookingGoApi/internal/models"

type CreateEmployeeRequest struct {
	UserID *int64 `json:"user_id"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

func (r *CreateEmployeeRequest) ToModel(userID, barbershopID int64) *models.Employee {
	return &models.Employee{
		UserID:       userID,
		BarberShopID: barbershopID,
	}
}
