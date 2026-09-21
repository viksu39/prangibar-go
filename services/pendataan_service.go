package services

import (
	"prangibar-go/config"
	"prangibar-go/models"
)

type PendataanService struct{}

func NewPendataanService() *PendataanService {
	return &PendataanService{}
}

func (s *PendataanService) FindByToken(token string) (*models.Perusahaan, error) {
	var perusahaan models.Perusahaan
	if err := config.DB.Preload("Pendataan").Where("token = ?", token).First(&perusahaan).Error; err != nil {
		return nil, err
	}
	return &perusahaan, nil
}

func (s *PendataanService) Submit(token string, pendataan *models.Pendataan) (*models.Pendataan, error) {
	perusahaan, err := s.FindByToken(token)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	if perusahaan.SkalaUsaha != nil && *perusahaan.SkalaUsaha == models.SkalaUMKM {
		return nil, ErrTokenNotFound
	}

	var existing models.Pendataan
	if err := config.DB.Where("perusahaan_id = ?", perusahaan.ID).First(&existing).Error; err == nil {
		pendataan.ID = existing.ID
		pendataan.PerusahaanID = perusahaan.ID
		if err := config.DB.Save(pendataan).Error; err != nil {
			return nil, err
		}
	} else {
		pendataan.PerusahaanID = perusahaan.ID
		if err := config.DB.Create(pendataan).Error; err != nil {
			return nil, err
		}
	}

	perusahaan.Status = models.StatusSelesai
	config.DB.Save(perusahaan)

	return pendataan, nil
}

func (s *PendataanService) FindAll() ([]models.Pendataan, error) {
	var pendataan []models.Pendataan
	if err := config.DB.Preload("Perusahaan").Order("submitted_at DESC").Find(&pendataan).Error; err != nil {
		return nil, err
	}
	return pendataan, nil
}

func (s *PendataanService) FindOne(id uint) (*models.Pendataan, error) {
	var pendataan models.Pendataan
	if err := config.DB.Preload("Perusahaan").First(&pendataan, id).Error; err != nil {
		return nil, err
	}
	return &pendataan, nil
}
