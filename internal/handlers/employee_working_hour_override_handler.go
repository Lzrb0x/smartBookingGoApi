package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-gonic/gin"
)

type EmployeeWorkingHourOverrideHandler struct {
	repo         *repositories.EmployeeWorkingHourOverrideRepository
	employeeRepo *repositories.EmployeeRepository
}

func NewEmployeeWorkingHourOverrideHandler(repo *repositories.EmployeeWorkingHourOverrideRepository, employeeRepo *repositories.EmployeeRepository) *EmployeeWorkingHourOverrideHandler {
	return &EmployeeWorkingHourOverrideHandler{repo: repo, employeeRepo: employeeRepo}
}

func (h *EmployeeWorkingHourOverrideHandler) GetAll(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	// Verify employee exists in barbershop
	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	overrides, err := h.repo.FindByEmployeeAndBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromOverrideModels(overrides))
}

func (h *EmployeeWorkingHourOverrideHandler) GetByID(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	overrideID, err := strconv.ParseInt(c.Param("overrideId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// Verify employee exists in barbershop
	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	override, err := h.repo.FindByID(c.Request.Context(), overrideID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "exceção de horário não encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify override belongs to this employee
	if override.EmployeeID != employeeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "exceção de horário não encontrada para este funcionário"})
		return
	}

	c.JSON(http.StatusOK, dtos.FromOverrideModel(override))
}

func (h *EmployeeWorkingHourOverrideHandler) Create(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	var req dtos.CreateEmployeeWorkingHourOverrideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify employee exists in barbershop
	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date inválido (formato: YYYY-MM-DD)"})
		return
	}

	// Check for duplicate date
	duplicate, err := h.repo.ExistsByEmployeeAndDate(c.Request.Context(), employeeID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if duplicate {
		c.JSON(http.StatusConflict, gin.H{"error": "já existe uma exceção de horário para esta data"})
		return
	}

	// Parse times if not a day off
	var startTime, endTime *time.Time
	if !req.IsDayOff {
		if req.StartTime == "" || req.EndTime == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_time e end_time são obrigatórios quando não é folga"})
			return
		}

		st, err := time.Parse("15:04:05", req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_time inválido (formato: HH:MM:SS)"})
			return
		}
		et, err := time.Parse("15:04:05", req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time inválido (formato: HH:MM:SS)"})
			return
		}

		if et.Before(st) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time não pode ser antes de start_time"})
			return
		}

		startTime = &st
		endTime = &et
	}

	override := req.ToModel(employeeID, date, startTime, endTime)
	if err := h.repo.Create(c.Request.Context(), override); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.FromOverrideModel(override))
}

func (h *EmployeeWorkingHourOverrideHandler) Update(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	overrideID, err := strconv.ParseInt(c.Param("overrideId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// Verify employee exists in barbershop
	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	// Verify override exists and belongs to this employee
	override, err := h.repo.FindByID(c.Request.Context(), overrideID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "exceção de horário não encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if override.EmployeeID != employeeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "exceção de horário não encontrada para este funcionário"})
		return
	}

	var req dtos.UpdateEmployeeWorkingHourOverrideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date inválido (formato: YYYY-MM-DD)"})
		return
	}

	// Check for duplicate date (excluding current)
	duplicate, err := h.repo.ExistsByEmployeeAndDateExcluding(c.Request.Context(), employeeID, date, overrideID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if duplicate {
		c.JSON(http.StatusConflict, gin.H{"error": "já existe uma exceção de horário para esta data"})
		return
	}

	// Parse times if not a day off
	var startTime, endTime *time.Time
	if !req.IsDayOff {
		if req.StartTime == "" || req.EndTime == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_time e end_time são obrigatórios quando não é folga"})
			return
		}

		st, err := time.Parse("15:04:05", req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_time inválido (formato: HH:MM:SS)"})
			return
		}
		et, err := time.Parse("15:04:05", req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time inválido (formato: HH:MM:SS)"})
			return
		}

		if et.Before(st) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time não pode ser antes de start_time"})
			return
		}

		startTime = &st
		endTime = &et
	}

	updated := req.ToModel(overrideID, employeeID, date, startTime, endTime)
	if err := h.repo.Update(c.Request.Context(), updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromOverrideModel(updated))
}

func (h *EmployeeWorkingHourOverrideHandler) Delete(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	overrideID, err := strconv.ParseInt(c.Param("overrideId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// Verify employee exists in barbershop
	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	// Verify override belongs to this employee
	override, err := h.repo.FindByID(c.Request.Context(), overrideID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "exceção de horário não encontrada"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if override.EmployeeID != employeeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "exceção de horário não encontrada para este funcionário"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), overrideID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
