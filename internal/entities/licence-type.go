package entities

import (
	"time"

	"gorm.io/gorm"
)

type LicenceType struct {
	gorm.Model
	Name        string    `gorm:"not null;unique"`
	Description string    `gorm:"not null"`
	MaxSeats    int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`

	TenantLicences []TenantLicence `json:"tenant_licences" gorm:"foreignKey:LicenceTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
