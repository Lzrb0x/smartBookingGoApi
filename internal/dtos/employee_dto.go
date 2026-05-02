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

type StaffBarbershopContextResponse struct {
	ID             int64  `json:"id"`
	BarbershopName string `json:"barbershop_name"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	OwnerID        int64  `json:"owner_id"`
	EmployeeID     *int64 `json:"employee_id"`
	Role           string `json:"role"`
}

type StaffContextResponse struct {
	IsOwner     bool                             `json:"is_owner"`
	IsEmployee  bool                             `json:"is_employee"`
	Barbershops []StaffBarbershopContextResponse `json:"barbershops"`
}
