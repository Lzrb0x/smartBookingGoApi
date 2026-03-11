package dtos

import "github.com/Lzrb0x/smartBookingGoApi/internal/models"

type CreateServiceRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
}

func (r *CreateServiceRequest) ToModel() *models.Service {
	return &models.Service{
		Name:        r.Name,
		Description: r.Description,
	}
}

type UpdateServiceRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
}

func (r *UpdateServiceRequest) ToModel() *models.Service {
	return &models.Service{
		Name:        r.Name,
		Description: r.Description,
	}
}
