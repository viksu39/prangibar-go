package services

import (
	"errors"
	"strconv"

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
	err := config.DB.Preload("BlokIII").Preload("BlokII").Where("email = ? OR nik = ?", email, nik).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (s *PendataanMahasiswaService) FindByID(id int32) (*models.PendataanMahasiswa, error) {
	var data models.PendataanMahasiswa
	err := config.DB.Preload("BlokIII").Preload("BlokII").First(&data, id).Error
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

	hash := data.PIN
	if hash == "" {
		// Fallback: registration PIN stored on mahasiswa table
		var m models.Mahasiswa
		if err := config.DB.Where("email = ? OR nik = ?", email, nik).First(&m).Error; err == nil {
			hash = m.PIN
		}
	}
	if hash == "" {
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin)); err != nil {
		return false, nil
	}
	return true, nil
}

// normalizeUsahaRows stamps sequential roster labels on Blok II children.
func normalizeUsahaRows(data *models.PendataanMahasiswa) {
	for i := range data.BlokII {
		if data.BlokII[i].NomorUsaha == "" {
			data.BlokII[i].NomorUsaha = strconv.Itoa(i + 1)
		}
	}
}

// saveAssociations persists parent (embedded fields) then fully replaces has-many
// children so autosave never accumulates duplicate rows.
//
// Empty incoming slices are treated carefully so incomplete payloads never wipe data:
//   - BlokIII empty → keep existing rows (minRows:1, empty is never a valid final state)
//   - BlokII empty + b1r14 unanswered/nil → keep existing rows
//   - BlokII empty + b1r14==1 → keep existing rows (household still reported as having a business)
//   - BlokII empty + b1r14!=1 → clear rows (respondent answered "Tidak memiliki usaha")
func saveAssociations(db *gorm.DB, data *models.PendataanMahasiswa, create bool) error {
	normalizeUsahaRows(data)

	var err error
	if create {
		err = db.Omit("BlokII", "BlokIII").Create(data).Error
	} else {
		err = db.Omit("BlokII", "BlokIII").Save(data).Error
	}
	if err != nil {
		return err
	}

	// Cleanly replace BlokII
	if len(data.BlokII) > 0 {
		var validUsaha []models.BlokIIBlock
		for _, u := range data.BlokII {
			if u.NamaUsaha != "" || u.NomorUsaha != "" {
				u.ID = 0
				u.PendataanMahasiswaID = data.ID
				validUsaha = append(validUsaha, u)
			}
		}
		if len(validUsaha) > 0 {
			db.Where("pendataan_mahasiswa_id = ?", data.ID).Delete(&models.BlokIIBlock{})
			if err := db.Create(&validUsaha).Error; err != nil {
				return err
			}
		}
	} else if data.BlokI.B1R14 != nil && *data.BlokI.B1R14 != 1 {
		db.Where("pendataan_mahasiswa_id = ?", data.ID).Delete(&models.BlokIIBlock{})
	}

	// Cleanly replace BlokIII (skip ghost rows with empty name and NIK)
	if len(data.BlokIII) > 0 {
		var validMembers []models.BlokIIIAnggotaKeluarga
		for _, m := range data.BlokIII {
			if m.NamaAnggota != "" || m.NIKAnggota != "" {
				m.ID = 0
				m.PendataanMahasiswaID = data.ID
				validMembers = append(validMembers, m)
			}
		}
		if len(validMembers) == 0 && data.BlokI.NamaKepalaKeluarga != "" {
			validMembers = append(validMembers, models.BlokIIIAnggotaKeluarga{
				PendataanMahasiswaID: data.ID,
				NomorUrut:            "1",
				NamaAnggota:          data.BlokI.NamaKepalaKeluarga,
				NIKAnggota:           data.BlokI.NIKKepalaKeluarga,
				HubunganKeluarga:     intPtr(1),
				KeberadaanAnggota:    intPtr(1),
				AlamatDomisili:       data.BlokI.AlamatSesuaiKK,
			})
		}
		if len(validMembers) > 0 {
			db.Where("pendataan_mahasiswa_id = ?", data.ID).Delete(&models.BlokIIIAnggotaKeluarga{})
			if err := db.Create(&validMembers).Error; err != nil {
				return err
			}
		}
	} else if data.BlokI.NamaKepalaKeluarga != "" {
		var count int64
		db.Model(&models.BlokIIIAnggotaKeluarga{}).Where("pendataan_mahasiswa_id = ?", data.ID).Count(&count)
		if count == 0 {
			m := models.BlokIIIAnggotaKeluarga{
				PendataanMahasiswaID: data.ID,
				NomorUrut:            "1",
				NamaAnggota:          data.BlokI.NamaKepalaKeluarga,
				NIKAnggota:           data.BlokI.NIKKepalaKeluarga,
				HubunganKeluarga:     intPtr(1),
				KeberadaanAnggota:    intPtr(1),
				AlamatDomisili:       data.BlokI.AlamatSesuaiKK,
			}
			db.Create(&m)
		}
	}
	return nil
}

func intPtr(i int) *int {
	return &i
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

		// Backfill PIN if a pre-existing draft has an empty bcrypt hash
		if data.PIN == "" {
			var m models.Mahasiswa
			if err := config.DB.Where("email = ? OR nik = ?", email, nik).First(&m).Error; err == nil && m.PIN != "" {
				data.PIN = m.PIN
			}
		}

		if err := saveAssociations(config.DB, data, false); err != nil {
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

	// Copy bcrypt PIN from registration so submit/verify can succeed
	var m models.Mahasiswa
	if err := config.DB.Where("email = ? OR nik = ?", email, nik).First(&m).Error; err == nil {
		data.PIN = m.PIN
	}

	if err := saveAssociations(config.DB, data, true); err != nil {
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

	if data.PIN == "" {
		var m models.Mahasiswa
		if err := config.DB.Where("email = ? OR nik = ?", email, nik).First(&m).Error; err == nil && m.PIN != "" {
			data.PIN = m.PIN
		}
	}

	if err := saveAssociations(config.DB, data, false); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *PendataanMahasiswaService) Delete(id int32) error {
	return config.DB.Delete(&models.PendataanMahasiswa{}, id).Error
}

type MahasiswaStatus struct {
	ID        int32  `json:"id"`
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
