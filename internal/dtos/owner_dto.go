package dtos

import "github.com/Lzrb0x/smartBookingGoApi/internal/models"

type CreateOwnerRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func (r *CreateOwnerRequest) ToModel() *models.Owner {
	return &models.Owner{
		UserID: r.UserID,
	}
}
