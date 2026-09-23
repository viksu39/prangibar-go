package services

import (
	"prangibar-go/config"
	"prangibar-go/models"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) FindByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := config.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *AdminService) FindAll() ([]models.Admin, error) {
	var admins []models.Admin
	if err := config.DB.Order("createdAt ASC").Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

func (s *AdminService) FindOne(id int32) (*models.Admin, error) {
	var admin models.Admin
	if err := config.DB.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *AdminService) Create(admin *models.Admin) (*models.Admin, error) {
	var existing models.Admin
	if err := config.DB.Where("email = ?", admin.Email).First(&existing).Error; err == nil {
		return nil, ErrEmailExists
	}
	if err := config.DB.Create(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *AdminService) Update(id int32, email, name, password string) (*models.Admin, error) {
	admin, err := s.FindOne(id)
	if err != nil {
		return nil, err
	}

	if email != "" && email != admin.Email {
		var existing models.Admin
		if err := config.DB.Where("email = ?", email).First(&existing).Error; err == nil && existing.ID != id {
			return nil, ErrEmailExists
		}
		admin.Email = email
	}
	if name != "" {
		admin.Name = name
	}
	if password != "" {
		admin.Password = password
	}

	if err := config.DB.Save(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *AdminService) Delete(id int32) (*models.Admin, error) {
	admin, err := s.FindOne(id)
	if err != nil {
		return nil, err
	}
	if err := config.DB.Delete(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}
