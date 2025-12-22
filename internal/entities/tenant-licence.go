package entities

import (
	"time"

	"gorm.io/gorm"
)

type TenantLicence struct {
	gorm.Model
	TenantID      uint       `json:"tenant_id"`
	LicenceTypeID uint       `json:"licence_type_id"`
	LicenceKey    string     `json:"licence_key"`
	UsedSeats     int        `json:"licenced_seats"`
	ExpiryDate    *time.Time `json:"expiry_date"`

	Tenant      Tenant      `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	LicenceType LicenceType `json:"licence_type" gorm:"foreignKey:LicenceTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
