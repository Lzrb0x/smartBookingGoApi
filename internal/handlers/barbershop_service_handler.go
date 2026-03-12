package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-gonic/gin"
)

type BarbershopServiceHandler struct {
	repo *repositories.BarbershopServiceRepository
}

func NewBarbershopServiceHandler(repo *repositories.BarbershopServiceRepository) *BarbershopServiceHandler {
	return &BarbershopServiceHandler{repo: repo}
}

func (h *BarbershopServiceHandler) GetAll(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	services, err := h.repo.FindByBarbershop(c.Request.Context(), barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func (h *BarbershopServiceHandler) GetByID(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	serviceID, err := strconv.ParseInt(c.Param("service_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do serviço inválido"})
		return
	}

	service, err := h.repo.FindByID(c.Request.Context(), barbershopID, serviceID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "serviço não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
}

func (h *BarbershopServiceHandler) Create(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	var req dtos.CreateBarbershopServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := req.ToModel(barbershopID)
	if err := h.repo.Create(c.Request.Context(), service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, service)
}

func (h *BarbershopServiceHandler) Update(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	serviceID, err := strconv.ParseInt(c.Param("service_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do serviço inválido"})
		return
	}

	var req dtos.UpdateBarbershopServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := req.ToModel(barbershopID)
	service.ID = serviceID

	if err := h.repo.Update(c.Request.Context(), service); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "serviço não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
}

func (h *BarbershopServiceHandler) Delete(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	serviceID, err := strconv.ParseInt(c.Param("service_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do serviço inválido"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), barbershopID, serviceID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "serviço não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
