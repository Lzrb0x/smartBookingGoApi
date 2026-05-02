package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	repo           *repositories.EmployeeRepository
	userRepo       *repositories.UserRepository
	ownerRepo      *repositories.OwnerRepository
	barbershopRepo *repositories.BarbershopRepository
}

func NewEmployeeHandler(
	repo *repositories.EmployeeRepository,
	userRepo *repositories.UserRepository,
	ownerRepo *repositories.OwnerRepository,
	barbershopRepo *repositories.BarbershopRepository,
) *EmployeeHandler {
	return &EmployeeHandler{
		repo:           repo,
		userRepo:       userRepo,
		ownerRepo:      ownerRepo,
		barbershopRepo: barbershopRepo,
	}
}

func (h *EmployeeHandler) GetAll(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	employees, err := h.repo.FindByBarbershop(c.Request.Context(), barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(employees) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "nenhum funcionário encontrado para esta barbearia"})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) GetStaffContext(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id inválido"})
		return
	}

	authenticatedUserID, ok := authenticatedUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}
	if authenticatedUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "você só pode consultar o próprio contexto profissional"})
		return
	}

	contextByBarbershop := make(map[int64]dtos.StaffBarbershopContextResponse)
	isOwner := false
	isEmployee := false

	owner, err := h.ownerRepo.FindByUserID(c.Request.Context(), userID)
	if err == nil {
		isOwner = true
		barbershops, err := h.barbershopRepo.FindByOwner(c.Request.Context(), owner.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, barbershop := range barbershops {
			contextByBarbershop[barbershop.ID] = dtos.StaffBarbershopContextResponse{
				ID:             barbershop.ID,
				BarbershopName: barbershop.BarbershopName,
				Address:        barbershop.Address,
				Phone:          barbershop.Phone,
				OwnerID:        barbershop.OwnerID,
				Role:           "owner",
			}
		}
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	employees, err := h.repo.FindByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, employee := range employees {
		isEmployee = true
		barbershop, err := h.barbershopRepo.FindByID(c.Request.Context(), employee.BarberShopID)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		employeeID := employee.ID
		context, exists := contextByBarbershop[barbershop.ID]
		if exists {
			context.EmployeeID = &employeeID
			context.Role = "owner_employee"
			contextByBarbershop[barbershop.ID] = context
			continue
		}

		contextByBarbershop[barbershop.ID] = dtos.StaffBarbershopContextResponse{
			ID:             barbershop.ID,
			BarbershopName: barbershop.BarbershopName,
			Address:        barbershop.Address,
			Phone:          barbershop.Phone,
			OwnerID:        barbershop.OwnerID,
			EmployeeID:     &employeeID,
			Role:           "employee",
		}
	}

	contexts := make([]dtos.StaffBarbershopContextResponse, 0, len(contextByBarbershop))
	for _, context := range contextByBarbershop {
		contexts = append(contexts, context)
	}
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].BarbershopName < contexts[j].BarbershopName
	})

	c.JSON(http.StatusOK, dtos.StaffContextResponse{
		IsOwner:     isOwner,
		IsEmployee:  isEmployee,
		Barbershops: contexts,
	})
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	var req dtos.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.resolveUserID(c, req)
	if err != nil {
		return
	}

	employee := req.ToModel(userID, barbershopID)
	if err := h.repo.Create(c.Request.Context(), employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) resolveUserID(c *gin.Context, req dtos.CreateEmployeeRequest) (int64, error) {
	var resolvedUser *models.User

	if req.UserID != nil {
		user, err := h.userRepo.FindByID(c.Request.Context(), *req.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado pelo id informado"})
				return 0, err
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 0, err
		}
		resolvedUser = user
	}

	email := strings.TrimSpace(req.Email)
	if email != "" {
		user, err := h.userRepo.FindByEmail(c.Request.Context(), email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado pelo email informado"})
				return 0, err
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 0, err
		}
		if resolvedUser != nil && resolvedUser.ID != user.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "os identificadores informados apontam para usuários diferentes"})
			return 0, errors.New("identificadores de usuários divergentes")
		}
		resolvedUser = user
	}

	phone := strings.TrimSpace(req.Phone)
	if phone != "" {
		user, err := h.userRepo.FindByPhone(c.Request.Context(), phone)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado pelo telefone informado"})
				return 0, err
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 0, err
		}
		if resolvedUser != nil && resolvedUser.ID != user.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "os identificadores informados apontam para usuários diferentes"})
			return 0, errors.New("identificadores de usuários divergentes")
		}
		resolvedUser = user
	}

	if resolvedUser == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "informe user_id, email e/ou phone para identificar o usuário"})
		return 0, errors.New("nenhum identificador de usuário informado")
	}

	return resolvedUser.ID, nil
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id da barbearia inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do funcionário inválido"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), barbershopID, employeeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
