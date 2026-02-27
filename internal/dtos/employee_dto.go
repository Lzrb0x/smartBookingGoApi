package dtos

import "github.com/Lzrb0x/smartBookingGoApi/internal/models"

type CreateEmployeeRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func (r *CreateEmployeeRequest) ToModel(barbershopID int64) *models.Employee {
	return &models.Employee{
		UserID:       r.UserID,
		BarberShopID: barbershopID,
	}
}
