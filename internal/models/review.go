package models

import (
	"time"
)

type ReviewStatus string

const (
	ReviewStatusUnderConsideration ReviewStatus = "under_consideration"
	ReviewStatusApproved           ReviewStatus = "approved"
	ReviewStatusRejected           ReviewStatus = "rejected"
)

type Review struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	MovieID      uint         `gorm:"not null;index" json:"movie_id"`
	UserID       uint         `gorm:"not null;index" json:"user_id"`
	Rating       int          `gorm:"not null;check:rating >= 1 AND rating <= 10" json:"rating"`
	Comment      string       `gorm:"type:text" json:"comment"`
	IsVisible    bool         `gorm:"default:true" json:"is_visible"` // админ может скрыть
	RejectReason string       `gorm:"type:text" json:"reject_reason"`
	Status       ReviewStatus `gorm:"default:0" json:"status"`
	Movie        Movie        `gorm:"foreignKey:MovieID" json:"-"`
	User         User         `gorm:"foreignKey:UserID" json:"user"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
