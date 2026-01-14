package models

type Genre struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null;unique" json:"name"`

	Movies []Movie `gorm:"many2many:movie_genres;" json:"movies,omitempty"`
}
