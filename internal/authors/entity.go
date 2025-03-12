package authors

import (
	"time"
	"toko-buku-api/internal/countries"
)

// Data models and structs specific to author functionality

type Author struct {
	ID         uint16
	Updated_At time.Time // `gorm:"updated_at"`
	Country_Id uint8
	Country    *countries.Country
	Author     string
	City       string
}

type CreateAuthorRequest struct {
	Country_Id uint8  `validate:"required" json:"country_id"`
	Author     string `validate:"required,min=3,max=50" json:"author"`
	City       string `validate:"required,min=3,max=50" json:"city"`
}

type UpdateAuthorRequest struct {
	ID         uint16 `validate:"required" json:"id"`
	Country_Id uint8  `validate:"required" json:"country_id"`
	Author     string `validate:"required,min=3,max=50" json:"author"`
	City       string `validate:"required,min=3,max=50" json:"city"`
}
