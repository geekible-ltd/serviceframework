package services

import (
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/entities"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
)

type LicenceTypeService struct {
	licenceTypeRepo *repositories.LicenceTypeRepository
}

func NewLicenceTypeService(licenceTypeRepo *repositories.LicenceTypeRepository) *LicenceTypeService {
	return &LicenceTypeService{licenceTypeRepo: licenceTypeRepo}
}

func (s *LicenceTypeService) GetAll() ([]frameworkdto.GetLicenceTypeDTO, error) {
	licences, err := s.licenceTypeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var licenceTypes []frameworkdto.GetLicenceTypeDTO
	for _, licence := range licences {
		licenceTypes = append(licenceTypes, frameworkdto.GetLicenceTypeDTO{
			ID:          licence.ID,
			Name:        licence.Name,
			Description: licence.Description,
			MaxSeats:    licence.MaxSeats,
		})
	}
	return licenceTypes, nil
}

func (s *LicenceTypeService) GetByID(id uint) (frameworkdto.GetLicenceTypeDTO, error) {
	licence, err := s.licenceTypeRepo.GetByID(id)
	if err != nil {
		return frameworkdto.GetLicenceTypeDTO{}, err
	}
	return frameworkdto.GetLicenceTypeDTO{
		ID:          licence.ID,
		Name:        licence.Name,
		Description: licence.Description,
		MaxSeats:    licence.MaxSeats,
	}, nil
}

func (s *LicenceTypeService) Create(dto frameworkdto.LicenceTypeCreateRequestDTO) error {
	licenceType := entities.LicenceType{
		Name:        dto.Name,
		Description: dto.Description,
		MaxSeats:    dto.MaxSeats,
	}
	if err := s.licenceTypeRepo.Create(licenceType, false); err != nil {
		return err
	}
	return nil
}

func (s *LicenceTypeService) Update(dto frameworkdto.LicenceTypeUpdateRequestDTO) error {
	licenceType, err := s.licenceTypeRepo.GetByID(dto.ID)
	if err != nil {
		return err
	}

	licenceType.Name = dto.Name
	licenceType.Description = dto.Description
	licenceType.MaxSeats = dto.MaxSeats

	if err := s.licenceTypeRepo.Update(licenceType); err != nil {
		return err
	}
	return nil
}

func (s *LicenceTypeService) Delete(id uint) error {
	licenceType, err := s.licenceTypeRepo.GetByID(id)
	if err != nil {
		return err
	}
	if err := s.licenceTypeRepo.Delete(licenceType); err != nil {
		return err
	}
	return nil
}
