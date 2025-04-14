package countries

import "time"

type Countries struct {
	ID           uint8
	Updated_At   time.Time
	Iso3         string
	Country      string
	Nice_Country string
	Currency     string
}

type CreateCountryRequest struct {
	Iso3         string `validate:"required,min=3,max=3" json:"iso3"`
	Country      string `validate:"required,min=3,max=50" json:"country"`
	Nice_Country string `validate:"required,min=3,max=50" json:"nice_country"`
	Currency     string `validate:"required,min=3,max=50" json:"currency"`
}

type UpdateCountryRequest struct {
	ID           uint8  `validate:"required" json:"id"`
	Iso3         string `validate:"required" json:"iso3"`
	Country      string `validate:"required,min=3,max=50" json:"country"`
	Nice_Country string `validate:"required,min=3,max=50" json:"nice_country"`
	Currency     string `validate:"required,min=3,max=50" json:"currency"`
}
