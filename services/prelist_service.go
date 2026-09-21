package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"prangibar-go/config"
	"prangibar-go/models"

	"github.com/xuri/excelize/v2"
)

type PrelistService struct{}

func NewPrelistService() *PrelistService {
	return &PrelistService{}
}

func (s *PrelistService) FindAll() ([]models.Perusahaan, error) {
	var perusahaan []models.Perusahaan
	if err := config.DB.Order("created_at DESC").Find(&perusahaan).Error; err != nil {
		return nil, err
	}
	return perusahaan, nil
}

func (s *PrelistService) FindOne(id uint) (*models.Perusahaan, error) {
	var perusahaan models.Perusahaan
	if err := config.DB.First(&perusahaan, id).Error; err != nil {
		return nil, err
	}
	return &perusahaan, nil
}

func (s *PrelistService) Create(perusahaan *models.Perusahaan) (*models.Perusahaan, error) {
	if err := config.DB.Create(perusahaan).Error; err != nil {
		return nil, err
	}
	return perusahaan, nil
}

func (s *PrelistService) Update(id uint, req interface{}) (*models.Perusahaan, error) {
	perusahaan, err := s.FindOne(id)
	if err != nil {
		return nil, err
	}

	updateReq, ok := req.(UpdatePrelistRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request type")
	}

	if updateReq.Nama != nil {
		perusahaan.Nama = *updateReq.Nama
	}
	if updateReq.Alamat != nil {
		perusahaan.Alamat = updateReq.Alamat
	}
	if updateReq.ContactPerson != nil {
		perusahaan.ContactPerson = updateReq.ContactPerson
	}
	if updateReq.Email != nil {
		perusahaan.Email = updateReq.Email
	}
	if updateReq.Phone != nil {
		perusahaan.Phone = updateReq.Phone
	}
	if updateReq.B1R1 != nil {
		perusahaan.B1R1 = updateReq.B1R1
	}
	if updateReq.B1R2 != nil {
		perusahaan.B1R2 = updateReq.B1R2
	}
	if updateReq.B1R3 != nil {
		perusahaan.B1R3 = updateReq.B1R3
	}
	if updateReq.B1R4 != nil {
		perusahaan.B1R4 = updateReq.B1R4
	}
	if updateReq.SkalaUsaha != nil {
		perusahaan.SkalaUsaha = updateReq.SkalaUsaha
	}

	if err := config.DB.Save(perusahaan).Error; err != nil {
		return nil, err
	}
	return perusahaan, nil
}

func (s *PrelistService) Delete(id uint) (*models.Perusahaan, error) {
	perusahaan, err := s.FindOne(id)
	if err != nil {
		return nil, err
	}
	if err := config.DB.Delete(perusahaan).Error; err != nil {
		return nil, err
	}
	return perusahaan, nil
}

func (s *PrelistService) UpdateToken(perusahaan *models.Perusahaan) error {
	return config.DB.Save(perusahaan).Error
}

func (s *PrelistService) BulkGenerateToken() (map[string]interface{}, error) {
	var perusahaan []models.Perusahaan
	if err := config.DB.Where("token IS NULL").Find(&perusahaan).Error; err != nil {
		return nil, err
	}

	type Result struct {
		ID    uint   `json:"id"`
		Nama  string `json:"nama"`
		Token string `json:"token"`
		URL   string `json:"url"`
	}

	var results []Result
	for _, p := range perusahaan {
		tokenBytes := make([]byte, 32)
		rand.Read(tokenBytes)
		token := hex.EncodeToString(tokenBytes)

		p.Token = &token
		config.DB.Save(&p)

		results = append(results, Result{
			ID:    p.ID,
			Nama:  p.Nama,
			Token: token,
			URL:   fmt.Sprintf("%s/form?token=%s", config.App.FrontendURL, token),
		})
	}

	return map[string]interface{}{
		"generated": len(results),
		"data":      results,
	}, nil
}

func (s *PrelistService) GenerateImportTemplate() []byte {
	f := excelize.NewFile()
	sheet := "Prelist"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"nama", "alamat", "contactPerson", "email", "phone", "kode_provinsi", "kode_kabupaten", "kode_kecamatan", "kode_desa"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	sample := []string{"PT. Contoh Sejahtera", "Jl. Sudirman No. 10, Medan", "Budi Santoso", "budi@contoh.co.id", "061-12345678", "12", "1271", "1271010", "1271010001"}
	for i, v := range sample {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, cell, v)
	}

	var buf bytes.Buffer
	f.Write(&buf)
	return buf.Bytes()
}

func (s *PrelistService) BulkImport(buffer []byte) (map[string]interface{}, error) {
	f, err := excelize.OpenReader(bytes.NewReader(buffer))
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file Excel")
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("file Excel tidak memiliki sheet")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet kosong atau tidak ada data")
	}

	header := rows[0]
	namaIdx := -1
	for i, h := range header {
		if h == "nama" {
			namaIdx = i
			break
		}
	}
	if namaIdx == -1 {
		return nil, fmt.Errorf("kolom wajib \"nama\" tidak ditemukan")
	}

	imported := 0
	skipped := []map[string]interface{}{}

	for i, row := range rows[1:] {
		nama := ""
		if namaIdx < len(row) {
			nama = row[namaIdx]
		}
		if nama == "" {
			skipped = append(skipped, map[string]interface{}{
				"row":    i + 2,
				"reason": "Kolom nama kosong",
			})
			continue
		}

		perusahaan := models.Perusahaan{Nama: nama}
		if len(row) > 1 {
			perusahaan.Alamat = &row[1]
		}
		if len(row) > 2 {
			perusahaan.ContactPerson = &row[2]
		}
		if len(row) > 3 {
			perusahaan.Email = &row[3]
		}
		if len(row) > 4 {
			perusahaan.Phone = &row[4]
		}
		if len(row) > 5 {
			perusahaan.B1R1 = &row[5]
		}
		if len(row) > 6 {
			perusahaan.B1R2 = &row[6]
		}
		if len(row) > 7 {
			perusahaan.B1R3 = &row[7]
		}
		if len(row) > 8 {
			perusahaan.B1R4 = &row[8]
		}

		config.DB.Create(&perusahaan)
		imported++
	}

	return map[string]interface{}{
		"imported":    imported,
		"skipped":     len(skipped),
		"skippedRows": skipped,
	}, nil
}

type UpdatePrelistRequest struct {
	Nama          *string
	Alamat        *string
	ContactPerson *string
	Email         *string
	Phone         *string
	B1R1          *string
	B1R2          *string
	B1R3          *string
	B1R4          *string
	SkalaUsaha    *models.SkalaUsaha
}
