package repositories

import (
	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	"github.com/geekible-ltd/serviceframework/internal/entities"
	"gorm.io/gorm"
)

type LicenceTypeRepository struct {
	db *gorm.DB
}

func NewLicenceTypeRepository(db *gorm.DB) *LicenceTypeRepository {
	return &LicenceTypeRepository{db: db}
}

func (r *LicenceTypeRepository) GetAll() ([]entities.LicenceType, error) {
	var licences []entities.LicenceType
	if err := r.db.Find(&licences).Error; err != nil {
		return nil, err
	}
	return licences, nil
}

func (r *LicenceTypeRepository) GetByID(id uint) (entities.LicenceType, error) {
	var licenceType entities.LicenceType
	if err := r.db.First(&licenceType, id).Error; err != nil {
		return entities.LicenceType{}, err
	}
	return licenceType, nil
}

func (r *LicenceTypeRepository) Create(licenceType entities.LicenceType, forSeeder bool) error {
	var licence entities.LicenceType
	if err := r.db.Where("name = ?", licenceType.Name).First(&licence).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if licence.ID != 0 {
		if !forSeeder {
			return frameworkconstants.ErrLicenceTypeAlreadyExists
		}
		return nil
	}

	return r.db.Create(&licenceType).Error
}

func (r *LicenceTypeRepository) Update(licenceType entities.LicenceType) error {
	return r.db.Save(&licenceType).Error
}

func (r *LicenceTypeRepository) Delete(licenceType entities.LicenceType) error {
	return r.db.Delete(&licenceType).Error
}
