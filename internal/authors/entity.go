package authors

import (
	"time"
	"toko-buku-api/internal/countries"
)

// Data models and structs specific to author functionality

type Authors struct {
	ID         uint16
	Updated_At time.Time // `gorm:"updated_at"`
	Country_Id uint8
	Country    *countries.Countries `json:",omitempty"`
	Author     string
	City       string
}

// func (a Author) TableName() string {
// 	return "authors"
// }

type CreateAuthorRequest struct {
	Country_Id uint8  `validate:"required" json:"country_id"`
	Author     string `validate:"required,min=3,max=50" json:"author"`
	City       string `validate:"required,min=3,max=50" json:"city"`
}

type UpdateAuthorRequest struct {
	ID         uint16 `json:"id"`         // `validate:"required"`
	Country_Id uint8  `json:"country_id"` // `validate:"required"`
	Author     string `json:"author"`     // `validate:"required,min=3,max=50"`
	City       string `json:"city"`       // `validate:"required,min=3,max=50"`
}
