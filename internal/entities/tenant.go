package entities

import (
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	IsActive bool   `json:"is_active"`

	Users         []User         `json:"users" gorm:"foreignKey:TenantID"`
	TenantLicence *TenantLicence `json:"tenant_licence,omitempty" gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
