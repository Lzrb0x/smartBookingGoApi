package models

import "time"

// EntityBase holds common fields shared by all entities.
type EntityBase struct {
	ID        int64     `db:"id"         json:"id"`
	Active    bool      `db:"active"     json:"active"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
}

// User identifies any user of the system (customer, owner, employee).
type User struct {
	EntityBase
	UserIdentifier string    `db:"user_identifier" json:"user_identifier"`
	Name           string    `db:"name"            json:"name"`
	Email          string    `db:"email"           json:"email"`
	Password       string    `db:"password"        json:"-"`
	Phone          string    `db:"phone"           json:"phone"`
	IsComplete     bool      `db:"is_complete"     json:"is_complete"`
	Bookings       []Booking `db:"-"               json:"bookings,omitempty"`
}

// Owner represents the owner of a barbershop.
type Owner struct {
	ID          int64        `db:"id"          json:"id"`
	UserID      int64        `db:"user_id"     json:"user_id"`
	User        *User        `db:"-"           json:"user,omitempty"`
	Barbershops []Barbershop `db:"-"           json:"barbershops,omitempty"`
}

// Barbershop represents a barbershop establishment.
type Barbershop struct {
	ID             int64               `db:"id"               json:"id"`
	BarbershopName string              `db:"barbershop_name"  json:"barbershop_name"`
	Address        string              `db:"address"          json:"address"`
	Phone          string              `db:"phone"            json:"phone"`
	OwnerID        int64               `db:"owner_id"         json:"owner_id"`
	Owner          *Owner              `db:"-"                json:"owner,omitempty"`
	Employees      []Employee          `db:"-"                json:"employees,omitempty"`
	Services       []BarbershopService `db:"-"         json:"services,omitempty"`
	Bookings       []Booking           `db:"-"                json:"bookings,omitempty"`
}

// Employee represents a barber/employee at a barbershop.
type Employee struct {
	ID                           int64                         `db:"id"             json:"id"`
	UserID                       int64                         `db:"user_id"        json:"user_id"`
	User                         *User                         `db:"-"              json:"user,omitempty"`
	BarberShopID                 int64                         `db:"barbershop_id"  json:"barbershop_id"`
	BarberShop                   *Barbershop                   `db:"-"              json:"barbershop,omitempty"`
	ServicesEmployee             []ServiceEmployee             `db:"-"              json:"services_employee,omitempty"`
	EmployeeWorkingHours         []EmployeeWorkingHour         `db:"-"              json:"working_hours,omitempty"`
	EmployeeWorkingHourOverrides []EmployeeWorkingHourOverride `db:"-"             json:"working_hour_overrides,omitempty"`
	Bookings                     []Booking                     `db:"-"              json:"bookings,omitempty"`
}

// Service is a globally catalogued service.
type Service struct {
	ID                 int64               `db:"id"          json:"id"`
	Name               string              `db:"name"        json:"name"`
	Description        string              `db:"description" json:"description"`
	BarbershopServices []BarbershopService `db:"-"          json:"barbershop_services,omitempty"`
}

// BarbershopService is a service customised for a specific barbershop.
type BarbershopService struct {
	ID                  int64             `db:"id"                   json:"id"`
	Price               float64           `db:"price"                json:"price"`
	Duration            int               `db:"duration"             json:"duration"`
	DescriptionOverride string            `db:"description_override" json:"description_override"`
	BarbershopID        int64             `db:"barbershop_id"        json:"barbershop_id"`
	ServiceID           int64             `db:"service_id"           json:"service_id"`
	Barbershop          *Barbershop       `db:"-"                    json:"barbershop,omitempty"`
	Service             *Service          `db:"-"                    json:"service,omitempty"`
	ServicesEmployees   []ServiceEmployee `db:"-"                    json:"services_employees,omitempty"`
}

// ServiceEmployee is the N-N association between an employee and a barbershop service.
type ServiceEmployee struct {
	ID                  int64 `db:"id"                    json:"id"`
	EmployeeID          int64 `db:"employee_id"           json:"employee_id"`
	BarbershopServiceID int64 `db:"barbershop_service_id" json:"barbershop_service_id"`
}

// EmployeeWorkingHour holds the standard weekly schedule for an employee.
type EmployeeWorkingHour struct {
	ID         int64     `db:"id"          json:"id"`
	EmployeeID int64     `db:"employee_id" json:"employee_id"`
	Employee   *Employee `db:"-"           json:"employee,omitempty"`
	DayOfWeek  int       `db:"day_of_week" json:"day_of_week"`
	StartTime  time.Time `db:"start_time"  json:"start_time"`
	EndTime    time.Time `db:"end_time"    json:"end_time"`
	IsDayOff   bool      `db:"is_day_off"  json:"is_day_off"`
}

// EmployeeWorkingHourOverride holds one-off schedule exceptions for an employee.
type EmployeeWorkingHourOverride struct {
	ID         int64     `db:"id"          json:"id"`
	EmployeeID int64     `db:"employee_id" json:"employee_id"`
	Employee   *Employee `db:"-"           json:"employee,omitempty"`
	Date       time.Time `db:"date"        json:"date"`
	StartTime  time.Time `db:"start_time"  json:"start_time"`
	EndTime    time.Time `db:"end_time"    json:"end_time"`
	IsDayOff   bool      `db:"is_day_off"  json:"is_day_off"`
}

// Booking is a confirmed appointment.
type Booking struct {
	ID                  int64              `db:"id"                    json:"id"`
	CustomerID          int64              `db:"customer_id"           json:"customer_id"`
	Customer            *User              `db:"-"                     json:"customer,omitempty"`
	EmployeeID          int64              `db:"employee_id"           json:"employee_id"`
	Employee            *Employee          `db:"-"                     json:"employee,omitempty"`
	BarbershopID        int64              `db:"barbershop_id"         json:"barbershop_id"`
	Barbershop          *Barbershop        `db:"-"                     json:"barbershop,omitempty"`
	BarbershopServiceID int64              `db:"barbershop_service_id" json:"barbershop_service_id"`
	BarbershopService   *BarbershopService `db:"-"                     json:"barbershop_service,omitempty"`
	Date                time.Time          `db:"date"                  json:"date"`
	StartTime           time.Time          `db:"start_time"            json:"start_time"`
	EndTime             time.Time          `db:"end_time"              json:"end_time"`
}
