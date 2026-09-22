package services

import (
	"errors"

	"prangibar-go/config"
	"prangibar-go/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PendataanMahasiswaService struct{}

func NewPendataanMahasiswaService() *PendataanMahasiswaService {
	return &PendataanMahasiswaService{}
}

func (s *PendataanMahasiswaService) FindByEmailOrNIK(email, nik string) (*models.PendataanMahasiswa, error) {
	var data models.PendataanMahasiswa
	err := config.DB.Preload("BlokIII").Where("email = ? OR nik = ?", email, nik).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (s *PendataanMahasiswaService) FindByID(id uint) (*models.PendataanMahasiswa, error) {
	var data models.PendataanMahasiswa
	err := config.DB.Preload("BlokIII").First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *PendataanMahasiswaService) VerifyPIN(email, nik, pin string) (bool, error) {
	var data models.PendataanMahasiswa
	err := config.DB.Where("email = ? AND nik = ?", email, nik).First(&data).Error
	if err != nil {
		return false, errors.New("Data tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.PIN), []byte(pin)); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *PendataanMahasiswaService) SaveDraft(email, nik string, data *models.PendataanMahasiswa) (*models.PendataanMahasiswa, error) {
	var existing models.PendataanMahasiswa
	err := config.DB.Where("email = ? OR nik = ?", email, nik).First(&existing).Error

	if err == nil {
		data.ID = existing.ID
		data.Email = existing.Email
		data.NIK = existing.NIK
		data.PIN = existing.PIN
		data.CreatedAt = existing.CreatedAt
		data.Status = "DRAFT"

		if err := config.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(data).Error; err != nil {
			return nil, err
		}
		return data, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	data.Email = email
	data.NIK = nik
	data.Status = "DRAFT"

	if err := config.DB.Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *PendataanMahasiswaService) Submit(email, nik, pin string, data *models.PendataanMahasiswa) (*models.PendataanMahasiswa, error) {
	valid, err := s.VerifyPIN(email, nik, pin)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("PIN salah")
	}

	var existing models.PendataanMahasiswa
	err = config.DB.Where("email = ? AND nik = ?", email, nik).First(&existing).Error
	if err != nil {
		return nil, errors.New("Data tidak ditemukan")
	}

	data.ID = existing.ID
	data.Email = existing.Email
	data.NIK = existing.NIK
	data.PIN = existing.PIN
	data.CreatedAt = existing.CreatedAt
	data.Status = "SUBMITTED"

	if err := config.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *PendataanMahasiswaService) Delete(id uint) error {
	return config.DB.Delete(&models.PendataanMahasiswa{}, id).Error
}

type MahasiswaStatus struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	NIK       string `json:"nik"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (s *PendataanMahasiswaService) ListWithStatus() ([]MahasiswaStatus, error) {
	var results []MahasiswaStatus
	err := config.DB.Model(&models.PendataanMahasiswa{}).
		Select("id, email, nik, status, created_at as createdAt, updated_at as updatedAt").
		Order("createdAt DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
