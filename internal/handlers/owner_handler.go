package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

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

func (h *OwnerHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id inválido"})
		return
	}

	owner, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "owner não encontrado para este usuário"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, owner)
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
