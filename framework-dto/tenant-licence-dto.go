package frameworkdto

import "time"

type GetTenantLicenceDTO struct {
	TenantLicenceID uint      `json:"tenant_licence_id"`
	TenantID        uint      `json:"tenant_id"`
	LicenceKey      string    `json:"licence_key"`
	LicenceType     string    `json:"licence_type"`
	LicenceStatus   string    `json:"licence_status"`
	LicenceExpiry   time.Time `json:"licence_expiry"`
}
