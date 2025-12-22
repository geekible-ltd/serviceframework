package services

import (
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
)

type TenantService struct {
	tenantRepo *repositories.TenantRepository
}

func NewTenantService(tenantRepo *repositories.TenantRepository) *TenantService {
	return &TenantService{tenantRepo: tenantRepo}
}

func (s *TenantService) GetTenantByID(tenantID uint) (frameworkdto.GetTenantDTO, error) {
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return frameworkdto.GetTenantDTO{}, err
	}

	return frameworkdto.GetTenantDTO{
		TenantID:      tenant.ID,
		TenantName:    tenant.Name,
		TenantEmail:   tenant.Email,
		TenantPhone:   tenant.Phone,
		TenantAddress: tenant.Address,
	}, nil
}

func (s *TenantService) GetAllTenants() ([]frameworkdto.GetTenantDTO, error) {
	tenants, err := s.tenantRepo.GetAll()
	if err != nil {
		return nil, err
	}

	tenantsDTO := make([]frameworkdto.GetTenantDTO, len(tenants))
	for i, tenant := range tenants {
		tenantsDTO[i] = frameworkdto.GetTenantDTO{
			TenantID:      tenant.ID,
			TenantName:    tenant.Name,
			TenantEmail:   tenant.Email,
			TenantPhone:   tenant.Phone,
			TenantAddress: tenant.Address,
		}
	}
	return tenantsDTO, nil
}

func (s *TenantService) UpdateTenant(tenantID uint, tenantDTO frameworkdto.UpdateTenantDTO) error {
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return err
	}

	tenant.Name = tenantDTO.TenantName
	tenant.Email = tenantDTO.TenantEmail
	tenant.Phone = tenantDTO.TenantPhone
	tenant.Address = tenantDTO.TenantAddress

	return s.tenantRepo.Update(tenant)
}

func (s *TenantService) DeleteTenant(tenantID uint) error {
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return err
	}

	return s.tenantRepo.Delete(tenant)
}
