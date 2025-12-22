package frameworkconstants

type UserRole string

const (
	UserRoleTenantAdmin UserRole = "tenant_admin"
	UserRoleTenantUser  UserRole = "tenant_user"
	UserRoleSuperAdmin  UserRole = "super_admin"
	UserRoleSuperUser   UserRole = "super_user"
)

const MaxFailedLoginAttempts = 3
const TokenKey = "token"
