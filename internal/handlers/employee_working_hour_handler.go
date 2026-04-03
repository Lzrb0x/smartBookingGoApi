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

type EmployeeWorkingHourHandler struct {
	repo         *repositories.EmployeeWorkingHourRepository
	employeeRepo *repositories.EmployeeRepository
}

func NewEmployeeWorkingHourHandler(repo *repositories.EmployeeWorkingHourRepository, employeeRepo *repositories.EmployeeRepository) *EmployeeWorkingHourHandler {
	return &EmployeeWorkingHourHandler{repo: repo, employeeRepo: employeeRepo}
}

func (h *EmployeeWorkingHourHandler) GetAll(c *gin.Context) {
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

	whs, err := h.repo.FindByEmployeeAndBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromWorkingHourModels(whs))
}

func (h *EmployeeWorkingHourHandler) GetByID(c *gin.Context) {
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

	whID, err := strconv.ParseInt(c.Param("workingHourId"), 10, 64)
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

	wh, err := h.repo.FindByID(c.Request.Context(), whID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "horário não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify the working hour belongs to this employee
	if wh.EmployeeID != employeeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "horário não encontrado para este funcionário"})
		return
	}

	c.JSON(http.StatusOK, dtos.FromWorkingHourModel(wh))
}

func (h *EmployeeWorkingHourHandler) Create(c *gin.Context) {
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

	var req dtos.CreateEmployeeWorkingHourRequest
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

	wh := req.ToModel(employeeID, startTime, endTime)
	if err := h.repo.Create(c.Request.Context(), wh); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.FromWorkingHourModel(wh))
}

func (h *EmployeeWorkingHourHandler) Update(c *gin.Context) {
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

	whID, err := strconv.ParseInt(c.Param("workingHourId"), 10, 64)
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

	// Verify working hour exists and belongs to this employee
	wh, err := h.repo.FindByID(c.Request.Context(), whID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "horário não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if wh.EmployeeID != employeeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "horário não encontrado para este funcionário"})
		return
	}

	var req dtos.UpdateEmployeeWorkingHourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	updated := req.ToModel(whID, employeeID, startTime, endTime)
	if err := h.repo.Update(c.Request.Context(), updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromWorkingHourModel(updated))
}

func (h *EmployeeWorkingHourHandler) Delete(c *gin.Context) {
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

	whID, err := strconv.ParseInt(c.Param("workingHourId"), 10, 64)
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

	// Verify working hour belongs to this employee
	wh, err := h.repo.FindByID(c.Request.Context(), whID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "horário não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if wh.EmployeeID != employeeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "horário não encontrado para este funcionário"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), whID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
