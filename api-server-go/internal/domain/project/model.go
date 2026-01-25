package project

import "time"

type Project struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	GitURL       string    `json:"gitURL"`
	SubDomain    string    `json:"subDomain"`
	CustomDomain *string   `json:"customDomain"`
	UserID       string    `json:"userId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
