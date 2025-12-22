package services

import (
	"strings"
	"time"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/entities"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegistrationService struct {
	userRepo          *repositories.UserRepository
	tenantRepo        *repositories.TenantRepository
	tenantLicenceRepo *repositories.TenantLicenceRepository
	licenceTypeRepo   *repositories.LicenceTypeRepository
}

func NewUserRegistrationService(
	userRepo *repositories.UserRepository,
	tenantRepo *repositories.TenantRepository,
	tenantLicenceRepo *repositories.TenantLicenceRepository,
	licenceTypeRepo *repositories.LicenceTypeRepository) *UserRegistrationService {
	return &UserRegistrationService{
		userRepo:          userRepo,
		tenantRepo:        tenantRepo,
		tenantLicenceRepo: tenantLicenceRepo,
		licenceTypeRepo:   licenceTypeRepo}
}

func (s *UserRegistrationService) RegisterTenant(tenantDTO frameworkdto.TenantRegistrationDTO) error {
	emailDomain := strings.Split(tenantDTO.Email, "@")[1]
	_, err := s.tenantRepo.GetByEmailDomain(emailDomain)

	if err == nil {
		return frameworkconstants.ErrTenantAlreadyExists
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	tenant := entities.Tenant{
		Name:     tenantDTO.Name,
		Email:    tenantDTO.Email,
		Phone:    tenantDTO.Phone,
		Address:  tenantDTO.Address,
		IsActive: true,
	}

	if err := s.tenantRepo.Create(&tenant); err != nil {
		return frameworkconstants.ErrFailedToCreateTenant
	}

	tenantLicence := entities.TenantLicence{
		TenantID:      tenant.ID,
		LicenceKey:    uuid.New().String(),
		LicenceTypeID: tenantDTO.LicenceTypeID,
		ExpiryDate:    nil,
	}
	if err := s.tenantLicenceRepo.Create(&tenantLicence); err != nil {
		return frameworkconstants.ErrFailedToCreateTenant
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(tenantDTO.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return frameworkconstants.ErrFailedToHashPassword
	}

	user := entities.User{
		TenantID:                        tenant.ID,
		FirstName:                       tenantDTO.User.FirstName,
		LastName:                        tenantDTO.User.LastName,
		Email:                           tenantDTO.User.Email,
		PasswordHash:                    string(passwordHash),
		FailedLoginAttempts:             0,
		IsActive:                        true,
		Role:                            string(frameworkconstants.UserRoleTenantAdmin),
		LastLoginAt:                     nil,
		LastLoginIP:                     "",
		ResetPasswordToken:              "",
		ResetPasswordTokenExpiresAt:     nil,
		IsEmailVerified:                 false,
		EmailVerificationToken:          uuid.New().String(),
		EmailVerificationTokenExpiresAt: nil,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return frameworkconstants.ErrFailedToCreateUser
	}

	return nil
}

func (s *UserRegistrationService) RegisterUser(tenantId uint, userDTO frameworkdto.UserRegistrationDTO) error {
	_, err := s.userRepo.GetByEmail(userDTO.Email)

	if err == nil {
		return frameworkconstants.ErrTenantAlreadyExists
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	tenantLicence, err := s.tenantLicenceRepo.GetByID(tenantId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return frameworkconstants.ErrTenantLicenceNotFound
	} else if err != nil {
		return err
	}

	licenceType, err := s.licenceTypeRepo.GetByID(tenantLicence.LicenceTypeID)
	if err != nil {
		return err
	}

	if tenantLicence.UsedSeats >= licenceType.MaxSeats {
		return frameworkconstants.ErrTenantLicenceExceeded
	}

	if tenantLicence.ExpiryDate != nil && tenantLicence.ExpiryDate.Before(time.Now()) {
		return frameworkconstants.ErrTenantLicenceExpired
	}

	tenantLicence.UsedSeats++
	if err := s.tenantLicenceRepo.Update(tenantLicence); err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return frameworkconstants.ErrFailedToHashPassword
	}

	user := entities.User{
		TenantID:                        tenantId,
		FirstName:                       userDTO.FirstName,
		LastName:                        userDTO.LastName,
		Email:                           userDTO.Email,
		PasswordHash:                    string(passwordHash),
		FailedLoginAttempts:             0,
		IsActive:                        true,
		Role:                            string(frameworkconstants.UserRoleTenantUser),
		LastLoginAt:                     nil,
		LastLoginIP:                     "",
		ResetPasswordToken:              "",
		ResetPasswordTokenExpiresAt:     nil,
		IsEmailVerified:                 false,
		EmailVerificationToken:          uuid.New().String(),
		EmailVerificationTokenExpiresAt: nil,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return frameworkconstants.ErrFailedToCreateUser
	}

	return nil
}

func (s *UserRegistrationService) DeleteUser(tenantId uint, userId uint) error {
	user, err := s.userRepo.GetByID(tenantId, userId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return frameworkconstants.ErrUserNotFound
	} else if err != nil {
		return err
	}

	s.userRepo.Delete(user)

	tenantLicence, err := s.tenantLicenceRepo.GetByTenantID(tenantId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return frameworkconstants.ErrTenantLicenceNotFound
	} else if err != nil {
		return err
	}

	tenantLicence.UsedSeats--
	if err := s.tenantLicenceRepo.Update(tenantLicence); err != nil {
		return err
	}

	return s.userRepo.Update(user)
}
