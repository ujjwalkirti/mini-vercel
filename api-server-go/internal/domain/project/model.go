package project

import "time"

type Project struct {
	ID           string
	Name         string
	GitURL       string
	SubDomain    string
	CustomDomain *string
	UserID       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
