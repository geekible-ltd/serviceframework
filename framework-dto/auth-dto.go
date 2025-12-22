package frameworkdto

type TokenDTO struct {
	Sub       string `json:"sub"`
	TenantID  uint   `json:"tenant_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
}
