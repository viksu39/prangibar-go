package controllers

import (
	"net/http"
	"strconv"

	"prangibar-go/middleware"
	"prangibar-go/models"
	"prangibar-go/services"

	"github.com/gin-gonic/gin"
)

type PendataanMahasiswaController struct {
	service *services.PendataanMahasiswaService
}

func NewPendataanMahasiswaController() *PendataanMahasiswaController {
	return &PendataanMahasiswaController{
		service: services.NewPendataanMahasiswaService(),
	}
}

// @Summary      Get pendataan mahasiswa by email or NIK
// @Tags         Pendataan Mahasiswa
// @Param        email  query     string  true  "Email"
// @Param        nik    query     string  true  "NIK"
// @Success      200    {object}  models.PendataanMahasiswa
// @Failure      404    {object}  map[string]string
// @Router       /pendataan-mahasiswa [get]
func (ctrl *PendataanMahasiswaController) GetPendataan(c *gin.Context) {
	email := c.Query("email")
	nik := c.Query("nik")

	if email == "" && nik == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau NIK harus diisi"})
		return
	}

	data, err := ctrl.service.FindByEmailOrNIK(email, nik)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary      Get pendataan mahasiswa by ID
// @Tags         Pendataan Mahasiswa
// @Param        id   path      int  true  "ID Pendataan"
// @Success      200  {object}  models.PendataanMahasiswa
// @Failure      404  {object}  map[string]string
// @Router       /pendataan-mahasiswa/{id} [get]
func (ctrl *PendataanMahasiswaController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := ctrl.service.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary      Verify PIN mahasiswa
// @Tags         Pendataan Mahasiswa
// @Param        request  body      models.VerifyPINRequest  true  "Email, NIK, PIN"
// @Success      200      {object}  models.VerifyPINResponse
// @Failure      401      {object}  map[string]string
// @Failure      423      {object}  map[string]interface{}
// @Router       /pendataan-mahasiswa/verify-pin [post]
func (ctrl *PendataanMahasiswaController) VerifyPIN(c *gin.Context) {
	var req models.VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lockKey := req.Email + ":" + req.NIK
	if locked, remaining := middleware.IsPinLocked(lockKey); locked {
		c.JSON(http.StatusLocked, gin.H{
			"error":             "Terlalu banyak percobaan PIN salah. Akun dikunci sementara.",
			"retry_after_seconds": int(remaining.Seconds()),
		})
		return
	}

	valid, err := ctrl.service.VerifyPIN(req.Email, req.NIK, req.PIN)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !valid {
		if locked := middleware.RecordPinFailure(lockKey); locked {
			c.JSON(http.StatusLocked, gin.H{
				"error":             "Terlalu banyak percobaan PIN salah. Akun dikunci selama 15 menit.",
				"retry_after_seconds": 900,
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "message": "PIN salah"})
		return
	}

	middleware.ResetPinAttempts(lockKey)
	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "PIN valid"})
}

// @Summary      Save draft pendataan mahasiswa
// @Tags         Pendataan Mahasiswa
// @Param        email    query     string                          true  "Email"
// @Param        nik      query     string                          true  "NIK"
// @Param        request  body      models.PendataanMahasiswa      true  "Data pendataan"
// @Success      200      {object}  models.PendataanMahasiswa
// @Failure      400      {object}  map[string]string
// @Router       /pendataan-mahasiswa/draft [post]
func (ctrl *PendataanMahasiswaController) SaveDraft(c *gin.Context) {
	email := c.Query("email")
	nik := c.Query("nik")

	if email == "" || nik == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan NIK harus diisi"})
		return
	}

	var data models.PendataanMahasiswa
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.service.SaveDraft(email, nik, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Submit pendataan mahasiswa
// @Tags         Pendataan Mahasiswa
// @Param        email    query     string                          true  "Email"
// @Param        nik      query     string                          true  "NIK"
// @Param        request  body      models.PendataanMahasiswa      true  "Data pendataan"
// @Success      200      {object}  models.PendataanMahasiswa
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /pendataan-mahasiswa/submit [post]
func (ctrl *PendataanMahasiswaController) Submit(c *gin.Context) {
	email := c.Query("email")
	nik := c.Query("nik")

	if email == "" || nik == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan NIK harus diisi"})
		return
	}

	var req struct {
		PIN string `json:"pin" binding:"required"`
		models.PendataanMahasiswa
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.service.Submit(email, nik, req.PIN, &req.PendataanMahasiswa)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Delete pendataan mahasiswa
// @Tags         Pendataan Mahasiswa
// @Security     BearerAuth
// @Param        id   path      int  true  "ID Pendataan"
// @Success      200  {object}  map[string]string
// @Router       /pendataan-mahasiswa/{id} [delete]
func (ctrl *PendataanMahasiswaController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

// @Summary      List semua mahasiswa dengan status pendataan
// @Tags         Pendataan Mahasiswa
// @Security     BearerAuth
// @Success      200  {array}   services.MahasiswaStatus
// @Router       /pendataan-mahasiswa/list [get]
func (ctrl *PendataanMahasiswaController) ListWithStatus(c *gin.Context) {
	results, err := ctrl.service.ListWithStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
