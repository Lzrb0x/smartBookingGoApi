package dtos

import (
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type CreateEmployeeWorkingHourOverrideRequest struct {
	Date      string `json:"date" binding:"required"` // YYYY-MM-DD format
	StartTime string `json:"start_time"`              // HH:MM:SS format, can be empty if is_day_off
	EndTime   string `json:"end_time"`                // HH:MM:SS format, can be empty if is_day_off
	IsDayOff  bool   `json:"is_day_off"`
}

type UpdateEmployeeWorkingHourOverrideRequest struct {
	Date      string `json:"date" binding:"required"` // YYYY-MM-DD format
	StartTime string `json:"start_time"`              // HH:MM:SS format, can be empty if is_day_off
	EndTime   string `json:"end_time"`                // HH:MM:SS format, can be empty if is_day_off
	IsDayOff  bool   `json:"is_day_off"`
}

type EmployeeWorkingHourOverrideResponse struct {
	ID         int64  `json:"id"`
	EmployeeID int64  `json:"employee_id"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	IsDayOff   bool   `json:"is_day_off"`
}

func (r *CreateEmployeeWorkingHourOverrideRequest) ToModel(employeeID int64, date time.Time, startTime, endTime *time.Time) *models.EmployeeWorkingHourOverride {
	return &models.EmployeeWorkingHourOverride{
		EmployeeID: employeeID,
		Date:       date,
		StartTime:  startTime,
		EndTime:    endTime,
		IsDayOff:   r.IsDayOff,
	}
}

func (r *UpdateEmployeeWorkingHourOverrideRequest) ToModel(id, employeeID int64, date time.Time, startTime, endTime *time.Time) *models.EmployeeWorkingHourOverride {
	return &models.EmployeeWorkingHourOverride{
		ID:         id,
		EmployeeID: employeeID,
		Date:       date,
		StartTime:  startTime,
		EndTime:    endTime,
		IsDayOff:   r.IsDayOff,
	}
}

func FromOverrideModel(override *models.EmployeeWorkingHourOverride) EmployeeWorkingHourOverrideResponse {
	startTimeStr := ""
	endTimeStr := ""
	if !override.IsDayOff && override.StartTime != nil && override.EndTime != nil {
		startTimeStr = override.StartTime.Format("15:04:05")
		endTimeStr = override.EndTime.Format("15:04:05")
	}
	return EmployeeWorkingHourOverrideResponse{
		ID:         override.ID,
		EmployeeID: override.EmployeeID,
		Date:       override.Date.Format("2006-01-02"),
		StartTime:  startTimeStr,
		EndTime:    endTimeStr,
		IsDayOff:   override.IsDayOff,
	}
}

func FromOverrideModels(overrides []models.EmployeeWorkingHourOverride) []EmployeeWorkingHourOverrideResponse {
	responses := make([]EmployeeWorkingHourOverrideResponse, len(overrides))
	for i, override := range overrides {
		responses[i] = FromOverrideModel(&override)
	}
	return responses
}
