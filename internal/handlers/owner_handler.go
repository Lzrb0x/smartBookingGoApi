package handlers

import (
	"net/http"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-gonic/gin"
)

type OwnerHandler struct {
	repo *repositories.OwnerRepository
}

func NewOwnerHandler(repo *repositories.OwnerRepository) *OwnerHandler {
	return &OwnerHandler{repo: repo}
}

func (h *OwnerHandler) Create(c *gin.Context) {
	var req dtos.CreateOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner := req.ToModel()
	if err := h.repo.Create(c.Request.Context(), owner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, owner)
}
