package dtos

import (
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type CreateBookingRequest struct {
	CustomerID          int64  `json:"customer_id" binding:"required"`
	EmployeeID          int64  `json:"employee_id" binding:"required"`
	BarbershopServiceID int64  `json:"barbershop_service_id" binding:"required"`
	Date                string `json:"date" binding:"required"`       // YYYY-MM-DD
	StartTime           string `json:"start_time" binding:"required"` // HH:MM:SS
}

type BookingResponse struct {
	ID                  int64  `json:"id"`
	CustomerID          int64  `json:"customer_id"`
	EmployeeID          int64  `json:"employee_id"`
	BarbershopID        int64  `json:"barbershop_id"`
	BarbershopServiceID int64  `json:"barbershop_service_id"`
	Date                string `json:"date"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
}

type AvailabilityResponse struct {
	Date                string   `json:"date"`
	EmployeeID          int64    `json:"employee_id"`
	BarbershopServiceID int64    `json:"barbershop_service_id"`
	ServiceDuration     int      `json:"service_duration"`
	Slots               []string `json:"slots"`
}

func FromBookingModel(booking *models.Booking) BookingResponse {
	return BookingResponse{
		ID:                  booking.ID,
		CustomerID:          booking.CustomerID,
		EmployeeID:          booking.EmployeeID,
		BarbershopID:        booking.BarbershopID,
		BarbershopServiceID: booking.BarbershopServiceID,
		Date:                booking.Date.Format("2006-01-02"),
		StartTime:           booking.StartTime.Format("15:04:05"),
		EndTime:             booking.EndTime.Format("15:04:05"),
	}
}
