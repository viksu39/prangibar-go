package services

import (
	"errors"

	"prangibar-go/config"
	"prangibar-go/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MahasiswaService struct{}

func NewMahasiswaService() *MahasiswaService {
	return &MahasiswaService{}
}

func (s *MahasiswaService) CheckPendataan(email, nik string) (*models.CheckPendataanResponse, error) {
	var mahasiswa models.Mahasiswa
	err := config.DB.Where("email = ? OR nik = ?", email, nik).First(&mahasiswa).Error

	if err == nil {
		return &models.CheckPendataanResponse{
			Registered: true,
			Email:      mahasiswa.Email,
			NIK:        mahasiswa.NIK,
			Message:    "Anda sudah terdaftar. Silakan lanjut ke penginputan pendataan.",
		}, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.CheckPendataanResponse{
			Registered: false,
			Email:      email,
			NIK:        nik,
			Message:    "Anda belum terdaftar. Silakan lakukan registrasi.",
		}, nil
	}

	return nil, err
}

func (s *MahasiswaService) Register(req *models.RegisterMahasiswaRequest) (*models.RegisterMahasiswaResponse, error) {
	if req.PIN != req.PINConfirm {
		return nil, errors.New("PIN dan konfirmasi PIN tidak cocok")
	}

	var existing models.Mahasiswa
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, errors.New("Email sudah terdaftar")
	}
	if err := config.DB.Where("nik = ?", req.NIK).First(&existing).Error; err == nil {
		return nil, errors.New("NIK sudah terdaftar")
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Gagal memproses PIN")
	}

	mahasiswa := models.Mahasiswa{
		Email: req.Email,
		NIK:   req.NIK,
		PIN:   string(hashedPIN),
	}

	if err := config.DB.Create(&mahasiswa).Error; err != nil {
		return nil, err
	}

	return &models.RegisterMahasiswaResponse{
		Success: true,
		Message: "Registrasi berhasil. Silakan lanjut ke penginputan pendataan.",
		Email:   mahasiswa.Email,
		NIK:     mahasiswa.NIK,
	}, nil
}

func (s *MahasiswaService) VerifyPIN(email, nik, pin string) (*models.VerifyPINResponse, error) {
	var mahasiswa models.Mahasiswa
	err := config.DB.Where("email = ? AND nik = ?", email, nik).First(&mahasiswa).Error
	if err != nil {
		return nil, errors.New("Data tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(mahasiswa.PIN), []byte(pin)); err != nil {
		return &models.VerifyPINResponse{
			Valid:   false,
			Message: "PIN salah",
		}, nil
	}

	return &models.VerifyPINResponse{
		Valid:   true,
		Message: "PIN valid",
	}, nil
}
