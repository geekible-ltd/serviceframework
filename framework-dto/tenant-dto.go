package frameworkdto

type GetTenantDTO struct {
	TenantID      uint   `json:"tenant_id"`
	TenantName    string `json:"tenant_name"`
	TenantEmail   string `json:"tenant_email"`
	TenantPhone   string `json:"tenant_phone"`
	TenantAddress string `json:"tenant_address"`
}

type UpdateTenantDTO struct {
	TenantName    string `json:"tenant_name"`
	TenantEmail   string `json:"tenant_email"`
	TenantPhone   string `json:"tenant_phone"`
	TenantAddress string `json:"tenant_address"`
}
