package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-gonic/gin"
)

type ServiceEmployeeHandler struct {
	repo *repositories.ServiceEmployeeRepository
}

func NewServiceEmployeeHandler(repo *repositories.ServiceEmployeeRepository) *ServiceEmployeeHandler {
	return &ServiceEmployeeHandler{repo: repo}
}

func (h *ServiceEmployeeHandler) AssignService(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	var req dtos.AssignServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify barbershop service exists for this barbershop
	serviceExists, err := h.repo.BarbershopServiceExists(c.Request.Context(), req.BarbershopServiceID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !serviceExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "serviço não encontrado para esta barbearia"})
		return
	}

	// Verify employee exists in this barbershop
	employeeExists, err := h.repo.EmployeeExistsInBarbershop(c.Request.Context(), req.EmployeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !employeeExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	// Check if already assigned
	alreadyAssigned, err := h.repo.IsEmployeeAssignedToService(c.Request.Context(), req.EmployeeID, req.BarbershopServiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if alreadyAssigned {
		c.JSON(http.StatusConflict, gin.H{"error": "funcionário já atribuído a este serviço"})
		return
	}

	serviceEmployee := req.ToModel()
	if err := h.repo.AssignServiceToEmployee(c.Request.Context(), serviceEmployee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.FromModel(serviceEmployee))
}

func (h *ServiceEmployeeHandler) UnassignService(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	var req dtos.UnassignServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify barbershop service exists for this barbershop
	serviceExists, err := h.repo.BarbershopServiceExists(c.Request.Context(), req.BarbershopServiceID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !serviceExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "serviço não encontrado para esta barbearia"})
		return
	}

	// Verify employee exists in this barbershop
	employeeExists, err := h.repo.EmployeeExistsInBarbershop(c.Request.Context(), req.EmployeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !employeeExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	// Check if assignment exists
	assigned, err := h.repo.IsEmployeeAssignedToService(c.Request.Context(), req.EmployeeID, req.BarbershopServiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !assigned {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não está atribuído a este serviço"})
		return
	}

	if err := h.repo.UnassignServiceFromEmployee(c.Request.Context(), req.EmployeeID, req.BarbershopServiceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ServiceEmployeeHandler) GetServicesByEmployee(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	serviceEmployees, err := h.repo.GetServicesByEmployee(c.Request.Context(), employeeID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromModels(serviceEmployees))
}

func (h *ServiceEmployeeHandler) GetEmployeesByService(c *gin.Context) {
	serviceID, err := strconv.ParseInt(c.Param("serviceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_id inválido"})
		return
	}

	serviceEmployees, err := h.repo.GetEmployeesByService(c.Request.Context(), serviceID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromModels(serviceEmployees))
}

func (h *ServiceEmployeeHandler) IsAssigned(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	serviceID, err := strconv.ParseInt(c.Param("serviceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_id inválido"})
		return
	}

	assigned, err := h.repo.IsEmployeeAssignedToService(c.Request.Context(), employeeID, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assigned": assigned})
}
