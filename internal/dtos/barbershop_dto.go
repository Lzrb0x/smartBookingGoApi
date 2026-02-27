package dtos

import "github.com/Lzrb0x/smartBookingGoApi/internal/models"

type CreateBarbershopRequest struct {
	BarbershopName string `json:"barbershop_name" binding:"required"`
	OwnerID        int64  `json:"owner_id"        binding:"required"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
}

func (r *CreateBarbershopRequest) ToModel() *models.Barbershop {
	return &models.Barbershop{
		BarbershopName: r.BarbershopName,
		OwnerID:        r.OwnerID,
		Address:        r.Address,
		Phone:          r.Phone,
	}
}

type UpdateBarbershopRequest struct {
	BarbershopName string `json:"barbershop_name" binding:"required"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
}

func (r *UpdateBarbershopRequest) ToModel() *models.Barbershop {
	return &models.Barbershop{
		BarbershopName: r.BarbershopName,
		Address:        r.Address,
		Phone:          r.Phone,
	}
}
