package dtos

import (
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
)

type CreateEmployeeWorkingHourRequest struct {
	DayOfWeek int    `json:"day_of_week" binding:"min=0,max=6"`
	StartTime string `json:"start_time"` // HH:MM:SS format, can be empty if is_day_off
	EndTime   string `json:"end_time"`   // HH:MM:SS format, can be empty if is_day_off
	IsDayOff  bool   `json:"is_day_off"`
}

type UpdateEmployeeWorkingHourRequest struct {
	DayOfWeek int    `json:"day_of_week" binding:"min=0,max=6"`
	StartTime string `json:"start_time"` // HH:MM:SS format, can be empty if is_day_off
	EndTime   string `json:"end_time"`   // HH:MM:SS format, can be empty if is_day_off
	IsDayOff  bool   `json:"is_day_off"`
}

type EmployeeWorkingHourResponse struct {
	ID         int64  `json:"id"`
	EmployeeID int64  `json:"employee_id"`
	DayOfWeek  int    `json:"day_of_week"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	IsDayOff   bool   `json:"is_day_off"`
}

func (r *CreateEmployeeWorkingHourRequest) ToModel(employeeID int64, startTime, endTime *time.Time) *models.EmployeeWorkingHour {
	return &models.EmployeeWorkingHour{
		EmployeeID: employeeID,
		DayOfWeek:  r.DayOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
		IsDayOff:   r.IsDayOff,
	}
}

func (r *UpdateEmployeeWorkingHourRequest) ToModel(id, employeeID int64, startTime, endTime *time.Time) *models.EmployeeWorkingHour {
	return &models.EmployeeWorkingHour{
		ID:         id,
		EmployeeID: employeeID,
		DayOfWeek:  r.DayOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
		IsDayOff:   r.IsDayOff,
	}
}

func FromWorkingHourModel(wh *models.EmployeeWorkingHour) EmployeeWorkingHourResponse {
	startTimeStr := ""
	endTimeStr := ""
	if !wh.IsDayOff && wh.StartTime != nil && wh.EndTime != nil {
		startTimeStr = wh.StartTime.Format("15:04:05")
		endTimeStr = wh.EndTime.Format("15:04:05")
	}
	return EmployeeWorkingHourResponse{
		ID:         wh.ID,
		EmployeeID: wh.EmployeeID,
		DayOfWeek:  wh.DayOfWeek,
		StartTime:  startTimeStr,
		EndTime:    endTimeStr,
		IsDayOff:   wh.IsDayOff,
	}
}

func FromWorkingHourModels(whs []models.EmployeeWorkingHour) []EmployeeWorkingHourResponse {
	responses := make([]EmployeeWorkingHourResponse, len(whs))
	for i, wh := range whs {
		responses[i] = FromWorkingHourModel(&wh)
	}
	return responses
}
