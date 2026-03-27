package dtos

import "github.com/Lzrb0x/smartBookingGoApi/internal/models"

type AssignServiceRequest struct {
	EmployeeID          int64 `json:"employee_id" binding:"required"`
	BarbershopServiceID int64 `json:"barbershop_service_id" binding:"required"`
}

func (r *AssignServiceRequest) ToModel() *models.ServiceEmployee {
	return &models.ServiceEmployee{
		EmployeeID:          r.EmployeeID,
		BarbershopServiceID: r.BarbershopServiceID,
	}
}

type UnassignServiceRequest struct {
	EmployeeID          int64 `json:"employee_id" binding:"required"`
	BarbershopServiceID int64 `json:"barbershop_service_id" binding:"required"`
}

type ServiceEmployeeResponse struct {
	ID                  int64 `json:"id"`
	EmployeeID          int64 `json:"employee_id"`
	BarbershopServiceID int64 `json:"barbershop_service_id"`
}

func FromModel(se *models.ServiceEmployee) ServiceEmployeeResponse {
	return ServiceEmployeeResponse{
		ID:                  se.ID,
		EmployeeID:          se.EmployeeID,
		BarbershopServiceID: se.BarbershopServiceID,
	}
}

func FromModels(ses []models.ServiceEmployee) []ServiceEmployeeResponse {
	responses := make([]ServiceEmployeeResponse, len(ses))
	for i, se := range ses {
		responses[i] = FromModel(&se)
	}
	return responses
}
