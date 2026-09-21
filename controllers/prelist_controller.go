package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"prangibar-go/config"
	"prangibar-go/models"
	"prangibar-go/services"

	"github.com/gin-gonic/gin"
)

type PrelistController struct {
	prelistService *services.PrelistService
}

func NewPrelistController() *PrelistController {
	return &PrelistController{
		prelistService: services.NewPrelistService(),
	}
}

type CreatePrelistRequest struct {
	Nama          string             `json:"nama" binding:"required" example:"PT. Maju Bersama"`
	Alamat        *string            `json:"alamat" example:"Jl. Sudirman No. 10, Medan"`
	ContactPerson *string            `json:"contactPerson" example:"Budi Santoso"`
	Email         *string            `json:"email" example:"budi@majubersama.co.id"`
	Phone         *string            `json:"phone" example:"061-12345678"`
	B1R1          *string            `json:"b1r1" example:"12"`
	B1R2          *string            `json:"b1r2" example:"1271"`
	B1R3          *string            `json:"b1r3" example:"1271010"`
	B1R4          *string            `json:"b1r4" example:"1271010001"`
	SkalaUsaha    *models.SkalaUsaha `json:"skalaUsaha" example:"UB"`
}

type UpdatePrelistRequest struct {
	Nama          *string            `json:"nama"`
	Alamat        *string            `json:"alamat"`
	ContactPerson *string            `json:"contactPerson"`
	Email         *string            `json:"email"`
	Phone         *string            `json:"phone"`
	B1R1          *string            `json:"b1r1"`
	B1R2          *string            `json:"b1r2"`
	B1R3          *string            `json:"b1r3"`
	B1R4          *string            `json:"b1r4"`
	SkalaUsaha    *models.SkalaUsaha `json:"skalaUsaha"`
}

// @Summary      List semua perusahaan target
// @Tags         Prelist
// @Security     BearerAuth
// @Success      200  {array}   models.Perusahaan
// @Router       /prelist [get]
func (ctrl *PrelistController) FindAll(c *gin.Context) {
	perusahaan, err := ctrl.prelistService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perusahaan)
}

// @Summary      Detail perusahaan
// @Tags         Prelist
// @Security     BearerAuth
// @Param        id   path      int  true  "Perusahaan ID"
// @Success      200  {object}  models.Perusahaan
// @Failure      404  {object}  map[string]string
// @Router       /prelist/{id} [get]
func (ctrl *PrelistController) FindOne(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	perusahaan, err := ctrl.prelistService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perusahaan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, perusahaan)
}

// @Summary      Tambah perusahaan target baru
// @Tags         Prelist
// @Security     BearerAuth
// @Param        request  body      CreatePrelistRequest  true  "Perusahaan data"
// @Success      201      {object}  models.Perusahaan
// @Router       /prelist [post]
func (ctrl *PrelistController) Create(c *gin.Context) {
	var req CreatePrelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	perusahaan := models.Perusahaan{
		Nama:          req.Nama,
		Alamat:        req.Alamat,
		ContactPerson: req.ContactPerson,
		Email:         req.Email,
		Phone:         req.Phone,
		B1R1:          req.B1R1,
		B1R2:          req.B1R2,
		B1R3:          req.B1R3,
		B1R4:          req.B1R4,
		SkalaUsaha:    req.SkalaUsaha,
	}

	created, err := ctrl.prelistService.Create(&perusahaan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// @Summary      Update data perusahaan
// @Tags         Prelist
// @Security     BearerAuth
// @Param        id       path      int                  true  "Perusahaan ID"
// @Param        request  body      UpdatePrelistRequest  true  "Perusahaan data"
// @Success      200      {object}  models.Perusahaan
// @Failure      404      {object}  map[string]string
// @Router       /prelist/{id} [put]
func (ctrl *PrelistController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdatePrelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := ctrl.prelistService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perusahaan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// @Summary      Hapus perusahaan
// @Tags         Prelist
// @Security     BearerAuth
// @Param        id   path      int  true  "Perusahaan ID"
// @Success      200  {object}  models.Perusahaan
// @Failure      404  {object}  map[string]string
// @Router       /prelist/{id} [delete]
func (ctrl *PrelistController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	deleted, err := ctrl.prelistService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perusahaan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, deleted)
}

// @Summary      Generate token akses pendataan
// @Tags         Prelist
// @Security     BearerAuth
// @Param        id   path      int  true  "Perusahaan ID"
// @Success      201  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /prelist/{id}/token [post]
func (ctrl *PrelistController) GenerateToken(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	perusahaan, err := ctrl.prelistService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perusahaan tidak ditemukan"})
		return
	}

	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	perusahaan.Token = &token
	if perusahaan.Status == models.StatusSelesai {
		perusahaan.Status = models.StatusBelum
	}

	if err := ctrl.prelistService.UpdateToken(perusahaan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := fmt.Sprintf("%s/form?token=%s", config.App.FrontendURL, token)
	c.JSON(http.StatusCreated, gin.H{"token": token, "url": url})
}

// @Summary      Bulk generate token
// @Tags         Prelist
// @Security     BearerAuth
// @Success      201  {object}  map[string]interface{}
// @Router       /prelist/bulk-token [post]
func (ctrl *PrelistController) BulkGenerateToken(c *gin.Context) {
	result, err := ctrl.prelistService.BulkGenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// @Summary      Download template Excel untuk import prelist
// @Tags         Prelist
// @Security     BearerAuth
// @Success      200  {file}    binary
// @Router       /prelist/import/template [get]
func (ctrl *PrelistController) DownloadTemplate(c *gin.Context) {
	buffer := ctrl.prelistService.GenerateImportTemplate()
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer)
}

// @Summary      Import perusahaan dari file Excel
// @Tags         Prelist
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Param        file  formData  file  true  "File .xlsx atau .xls"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /prelist/import [post]
func (ctrl *PrelistController) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer f.Close()

	buffer := make([]byte, file.Size)
	f.Read(buffer)

	result, err := ctrl.prelistService.BulkImport(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
