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

type BookingDashboardItemResponse struct {
	ID                         int64   `json:"id"`
	CustomerID                 int64   `json:"customer_id"`
	EmployeeID                 int64   `json:"employee_id"`
	EmployeeName               string  `json:"employee_name"`
	BarbershopID               int64   `json:"barbershop_id"`
	BarbershopName             string  `json:"barbershop_name"`
	BarbershopAddress          string  `json:"barbershop_address"`
	BarbershopPhone            string  `json:"barbershop_phone"`
	BarbershopServiceID        int64   `json:"barbershop_service_id"`
	ServiceID                  int64   `json:"service_id"`
	ServiceName                string  `json:"service_name"`
	ServiceDescription         string  `json:"service_description"`
	ServicePrice               float64 `json:"service_price"`
	ServiceDuration            int     `json:"service_duration"`
	ServiceDescriptionOverride string  `json:"service_description_override"`
	Date                       string  `json:"date"`
	StartTime                  string  `json:"start_time"`
	EndTime                    string  `json:"end_time"`
}

type RecentBarbershopResponse struct {
	ID              int64  `json:"id"`
	BarbershopName  string `json:"barbershop_name"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	LastBookingDate string `json:"last_booking_date"`
	LastStartTime   string `json:"last_start_time"`
}

type UserDashboardResponse struct {
	CurrentBooking    *BookingDashboardItemResponse  `json:"current_booking"`
	RecentBookings    []BookingDashboardItemResponse `json:"recent_bookings"`
	RecentBarbershops []RecentBarbershopResponse     `json:"recent_barbershops"`
}

type ProfessionalBookingItemResponse struct {
	ID                  int64   `json:"id"`
	CustomerID          int64   `json:"customer_id"`
	CustomerName        string  `json:"customer_name"`
	CustomerPhone       string  `json:"customer_phone"`
	EmployeeID          int64   `json:"employee_id"`
	EmployeeUserID      int64   `json:"employee_user_id"`
	EmployeeName        string  `json:"employee_name"`
	EmployeePhone       string  `json:"employee_phone"`
	BarbershopID        int64   `json:"barbershop_id"`
	BarbershopName      string  `json:"barbershop_name"`
	BarbershopServiceID int64   `json:"barbershop_service_id"`
	ServiceID           int64   `json:"service_id"`
	ServiceName         string  `json:"service_name"`
	ServicePrice        float64 `json:"service_price"`
	ServiceDuration     int     `json:"service_duration"`
	Date                string  `json:"date"`
	StartTime           string  `json:"start_time"`
	EndTime             string  `json:"end_time"`
}

type ProfessionalBookingDashboardResponse struct {
	Date             string                            `json:"date"`
	EmployeeID       *int64                            `json:"employee_id"`
	EstimatedRevenue float64                           `json:"estimated_revenue"`
	AppointmentCount int                               `json:"appointment_count"`
	TotalDuration    int                               `json:"total_duration"`
	Bookings         []ProfessionalBookingItemResponse `json:"bookings"`
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

func FromBookingDashboardItemModel(booking *models.BookingDashboardItem) BookingDashboardItemResponse {
	return BookingDashboardItemResponse{
		ID:                         booking.ID,
		CustomerID:                 booking.CustomerID,
		EmployeeID:                 booking.EmployeeID,
		EmployeeName:               booking.EmployeeName,
		BarbershopID:               booking.BarbershopID,
		BarbershopName:             booking.BarbershopName,
		BarbershopAddress:          booking.BarbershopAddress,
		BarbershopPhone:            booking.BarbershopPhone,
		BarbershopServiceID:        booking.BarbershopServiceID,
		ServiceID:                  booking.ServiceID,
		ServiceName:                booking.ServiceName,
		ServiceDescription:         booking.ServiceDescription,
		ServicePrice:               booking.ServicePrice,
		ServiceDuration:            booking.ServiceDuration,
		ServiceDescriptionOverride: booking.ServiceDescriptionOverride,
		Date:                       booking.Date.Format("2006-01-02"),
		StartTime:                  booking.StartTime.Format("15:04:05"),
		EndTime:                    booking.EndTime.Format("15:04:05"),
	}
}

func FromBookingDashboardItemModels(bookings []models.BookingDashboardItem) []BookingDashboardItemResponse {
	responses := make([]BookingDashboardItemResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = FromBookingDashboardItemModel(&booking)
	}
	return responses
}

func FromProfessionalBookingItemModel(booking *models.ProfessionalBookingItem) ProfessionalBookingItemResponse {
	return ProfessionalBookingItemResponse{
		ID:                  booking.ID,
		CustomerID:          booking.CustomerID,
		CustomerName:        booking.CustomerName,
		CustomerPhone:       booking.CustomerPhone,
		EmployeeID:          booking.EmployeeID,
		EmployeeUserID:      booking.EmployeeUserID,
		EmployeeName:        booking.EmployeeName,
		EmployeePhone:       booking.EmployeePhone,
		BarbershopID:        booking.BarbershopID,
		BarbershopName:      booking.BarbershopName,
		BarbershopServiceID: booking.BarbershopServiceID,
		ServiceID:           booking.ServiceID,
		ServiceName:         booking.ServiceName,
		ServicePrice:        booking.ServicePrice,
		ServiceDuration:     booking.ServiceDuration,
		Date:                booking.Date.Format("2006-01-02"),
		StartTime:           booking.StartTime.Format("15:04:05"),
		EndTime:             booking.EndTime.Format("15:04:05"),
	}
}

func FromProfessionalBookingItemModels(bookings []models.ProfessionalBookingItem) []ProfessionalBookingItemResponse {
	responses := make([]ProfessionalBookingItemResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = FromProfessionalBookingItemModel(&booking)
	}
	return responses
}

func FromRecentBarbershopModel(barbershop *models.RecentBarbershop) RecentBarbershopResponse {
	return RecentBarbershopResponse{
		ID:              barbershop.ID,
		BarbershopName:  barbershop.BarbershopName,
		Address:         barbershop.Address,
		Phone:           barbershop.Phone,
		LastBookingDate: barbershop.LastBookingDate.Format("2006-01-02"),
		LastStartTime:   barbershop.LastStartTime.Format("15:04:05"),
	}
}

func FromRecentBarbershopModels(barbershops []models.RecentBarbershop) []RecentBarbershopResponse {
	responses := make([]RecentBarbershopResponse, len(barbershops))
	for i, barbershop := range barbershops {
		responses[i] = FromRecentBarbershopModel(&barbershop)
	}
	return responses
}
