package dtos

import "github.com/Lzrb0x/smartBookingGoApi/internal/models"

type CreateBarbershopServiceRequest struct {
	Price               float64 `json:"price"                 binding:"required"`
	Duration            int     `json:"duration"              binding:"required"`
	DescriptionOverride string  `json:"description_override"`
	ServiceID           int64   `json:"service_id"            binding:"required"`
}

func (r *CreateBarbershopServiceRequest) ToModel(barbershopID int64) *models.BarbershopService {
	return &models.BarbershopService{
		Price:               r.Price,
		Duration:            r.Duration,
		DescriptionOverride: r.DescriptionOverride,
		BarbershopID:        barbershopID,
		ServiceID:           r.ServiceID,
	}
}

type UpdateBarbershopServiceRequest struct {
	Price               float64 `json:"price"                 binding:"required"`
	Duration            int     `json:"duration"              binding:"required"`
	DescriptionOverride string  `json:"description_override"`
	ServiceID           int64   `json:"service_id"            binding:"required"`
}

func (r *UpdateBarbershopServiceRequest) ToModel(barbershopID int64) *models.BarbershopService {
	return &models.BarbershopService{
		Price:               r.Price,
		Duration:            r.Duration,
		DescriptionOverride: r.DescriptionOverride,
		BarbershopID:        barbershopID,
		ServiceID:           r.ServiceID,
	}
}
