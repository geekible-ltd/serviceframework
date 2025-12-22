package frameworkconstants

import "errors"

const (
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeValidation          = "VALIDATION_ERROR"
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabase            = "DATABASE_ERROR"
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeMissingHeader       = "MISSING_HEADER"
	ErrCodeInvalidUUID         = "INVALID_UUID"
	ErrCodeDuplicateEntry      = "DUPLICATE_ENTRY"
	ErrCodeForeignKeyViolation = "FOREIGN_KEY_VIOLATION"
	ErrCodeInvalidBody         = "INVALID_BODY"
	ErrUserAccountLocked       = "ACCOUNT_LOCKED"
	ErrUnauthorizedError       = "UNAUTHORIZED_ERROR"
)

var (
	ErrFailedToCreateTenant        = errors.New("failed to create tenant")
	ErrFailedToCreateUser          = errors.New("failed to create user")
	ErrFailedToHashPassword        = errors.New("failed to hash password")
	ErrTenantAlreadyExists         = errors.New("tenant already exists")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrUserNotFound                = errors.New("user not found")
	ErrInvalidPassword             = errors.New("invalid password")
	ErrTenantNotFound              = errors.New("tenant not found")
	ErrTenantLicenceNotFound       = errors.New("tenant licence not found")
	ErrTenantLicenceExceeded       = errors.New("tenant licence exceeded")
	ErrTenantLicenceExpired        = errors.New("tenant licence expired")
	ErrFailedToCreateTenantLicence = errors.New("failed to create tenant licence")
	ErrLicenceTypeAlreadyExists    = errors.New("licence type already exists")
)
