package entities

import (
	"time"

	"gorm.io/gorm"
)

type TenantLicence struct {
	gorm.Model
	TenantID      uint       `json:"tenant_id"`
	LicenceKey    string     `json:"licence_key"`
	LicencedSeats int        `json:"licenced_seats"`
	UsedSeats     int        `json:"used_seats"`
	ExpiryDate    *time.Time `json:"expiry_date"`

	Tenant Tenant `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
