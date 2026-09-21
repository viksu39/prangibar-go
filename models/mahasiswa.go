package models

import (
	"time"

	"gorm.io/gorm"
)

type Mahasiswa struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Email     string         `gorm:"unique;not null;index" json:"email"`
	NIK       string         `gorm:"unique;not null;index" json:"nik"`
	PIN       string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CheckPendataanRequest struct {
	Email string `json:"email" binding:"required,email" example:"mahasiswa@example.com"`
	NIK   string `json:"nik" binding:"required,len=16" example:"1234567890123456"`
}

type CheckPendataanResponse struct {
	Registered bool   `json:"registered" example:"true"`
	Email      string `json:"email" example:"mahasiswa@example.com"`
	NIK        string `json:"nik" example:"1234567890123456"`
	Message    string `json:"message" example:"Anda sudah terdaftar"`
}

type RegisterMahasiswaRequest struct {
	Email        string `json:"email" binding:"required,email" example:"mahasiswa@example.com"`
	NIK          string `json:"nik" binding:"required,len=16" example:"1234567890123456"`
	PIN          string `json:"pin" binding:"required,min=6,max=6" example:"123456"`
	PINConfirm   string `json:"pinConfirm" binding:"required,min=6,max=6" example:"123456"`
}

type RegisterMahasiswaResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Registrasi berhasil"`
	Email   string `json:"email" example:"mahasiswa@example.com"`
	NIK     string `json:"nik" example:"1234567890123456"`
}

type VerifyPINRequest struct {
	Email string `json:"email" binding:"required,email" example:"mahasiswa@example.com"`
	NIK   string `json:"nik" binding:"required,len=16" example:"1234567890123456"`
	PIN   string `json:"pin" binding:"required,min=6,max=6" example:"123456"`
}

type VerifyPINResponse struct {
	Valid   bool   `json:"valid" example:"true"`
	Message string `json:"message" example:"PIN valid"`
}
