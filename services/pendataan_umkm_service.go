package services

import (
	"prangibar-go/config"
	"prangibar-go/models"
)

type PendataanUmkmService struct{}

func NewPendataanUmkmService() *PendataanUmkmService {
	return &PendataanUmkmService{}
}

func (s *PendataanUmkmService) FindByToken(token string) (*models.Perusahaan, error) {
	var perusahaan models.Perusahaan
	if err := config.DB.Preload("PendataanUmkm").Where("token = ?", token).First(&perusahaan).Error; err != nil {
		return nil, ErrTokenNotFound
	}

	if perusahaan.SkalaUsaha == nil || *perusahaan.SkalaUsaha != models.SkalaUMKM {
		return nil, ErrTokenNotFound
	}

	return &perusahaan, nil
}

func (s *PendataanUmkmService) Submit(token string, pendataan *models.PendataanUmkm) (*models.PendataanUmkm, error) {
	perusahaan, err := s.FindByToken(token)
	if err != nil {
		return nil, err
	}

	var existing models.PendataanUmkm
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

func (s *PendataanUmkmService) FindAll() ([]models.PendataanUmkm, error) {
	var pendataan []models.PendataanUmkm
	if err := config.DB.Preload("Perusahaan").Order("submittedAt DESC").Find(&pendataan).Error; err != nil {
		return nil, err
	}
	return pendataan, nil
}

func (s *PendataanUmkmService) FindOne(id int32) (*models.PendataanUmkm, error) {
	var pendataan models.PendataanUmkm
	if err := config.DB.Preload("Perusahaan").First(&pendataan, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &pendataan, nil
}
