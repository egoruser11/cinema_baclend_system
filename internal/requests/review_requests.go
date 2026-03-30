package requests

import "cinema_backend_system/internal/models"

type ReviewCreateRequest struct {
	MovieID uint   `json:"movie_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment" binding:"required,min=1,max=20"`
}
type ReviewApproveRequest struct {
	ReviewID     uint                `json:"review_id" binding:"required"`
	Status       models.ReviewStatus `json:"status" binding:"required"`
	ReasonReject *string             `json:"reason" binding:"required,min=1,max=20"`
}
