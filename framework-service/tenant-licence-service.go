package frameworkservice

import (
	"time"

	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
)

type TenantLicenceService struct {
	tenantLicenceRepo *repositories.TenantLicenceRepository
	licenceTypeRepo   *repositories.LicenceTypeRepository
}

func NewTenantLicenceService(tenantLicenceRepo *repositories.TenantLicenceRepository, licenceTypeRepo *repositories.LicenceTypeRepository) *TenantLicenceService {
	return &TenantLicenceService{tenantLicenceRepo: tenantLicenceRepo, licenceTypeRepo: licenceTypeRepo}
}

func (s *TenantLicenceService) GetTenantLicenceByID(tenantID uint) (frameworkdto.GetTenantLicenceDTO, error) {
	licence, err := s.tenantLicenceRepo.GetByTenantID(tenantID)
	if err != nil {
		return frameworkdto.GetTenantLicenceDTO{}, err
	}

	licenceType, err := s.licenceTypeRepo.GetByID(licence.LicenceTypeID)
	if err != nil {
		return frameworkdto.GetTenantLicenceDTO{}, err
	}

	return frameworkdto.GetTenantLicenceDTO{
		TenantLicenceID: licence.ID,
		TenantID:        licence.TenantID,
		LicenceKey:      licence.LicenceKey,
		LicenceType:     licenceType.Name,
		LicenceExpiry:   *licence.ExpiryDate,
	}, nil
}

func (s *TenantLicenceService) GetTenantLicencesByTenantID(tenantID uint) (frameworkdto.GetTenantLicenceDTO, error) {
	licence, err := s.tenantLicenceRepo.GetByTenantID(tenantID)
	if err != nil {
		return frameworkdto.GetTenantLicenceDTO{}, err
	}

	licenceType, err := s.licenceTypeRepo.GetByID(licence.LicenceTypeID)
	if err != nil {
		return frameworkdto.GetTenantLicenceDTO{}, err
	}

	return frameworkdto.GetTenantLicenceDTO{
		TenantLicenceID: licence.ID,
		TenantID:        licence.TenantID,
		LicenceKey:      licence.LicenceKey,
		LicenceType:     licenceType.Name,
		LicenceExpiry:   *licence.ExpiryDate,
	}, nil
}

func (s *TenantLicenceService) RenewTenantLicence(tenantID uint, addDays int) error {
	licence, err := s.tenantLicenceRepo.GetByTenantID(tenantID)
	if err != nil {
		return err
	}

	expiry := time.Now().AddDate(0, 0, addDays)
	licence.ExpiryDate = &expiry

	err = s.tenantLicenceRepo.Update(licence)
	if err != nil {
		return err
	}

	return nil
}

func (s *TenantLicenceService) ChangeTenantLicenceType(tenantID, downgradeToLicenceTypeID uint, addDays int) error {
	licence, err := s.tenantLicenceRepo.GetByTenantID(tenantID)
	if err != nil {
		return err
	}

	if addDays > 0 {
		expiry := time.Now().AddDate(0, 0, addDays)
		licence.ExpiryDate = &expiry
	} else {
		licence.ExpiryDate = nil
	}

	licence.LicenceTypeID = downgradeToLicenceTypeID

	err = s.tenantLicenceRepo.Update(licence)
	if err != nil {
		return err
	}

	return nil
}
