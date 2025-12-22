package frameworkservice

import (
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
)

type TenantLicenceService struct {
	tenantLicenceRepo *repositories.TenantLicenceRepository
}

func NewTenantLicenceService(tenantLicenceRepo *repositories.TenantLicenceRepository) *TenantLicenceService {
	return &TenantLicenceService{tenantLicenceRepo: tenantLicenceRepo}
}

func (s *TenantLicenceService) GetTenantLicenceByID(tenantID uint) (frameworkdto.GetTenantLicenceDTO, error) {
	licence, err := s.tenantLicenceRepo.GetByTenantID(tenantID)
	if err != nil {
		return frameworkdto.GetTenantLicenceDTO{}, err
	}
	return frameworkdto.GetTenantLicenceDTO{
		TenantLicenceID: licence.ID,
		TenantID:        licence.TenantID,
		LicenceKey:      licence.LicenceKey,
		LicenceType:     "",
		LicenceExpiry:   *licence.ExpiryDate,
	}, nil
}
