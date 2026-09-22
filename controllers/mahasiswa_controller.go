package controllers

import (
	"net/http"

	"prangibar-go/middleware"
	"prangibar-go/models"
	"prangibar-go/services"

	"github.com/gin-gonic/gin"
)

type MahasiswaController struct {
	mahasiswaService *services.MahasiswaService
}

func NewMahasiswaController() *MahasiswaController {
	return &MahasiswaController{
		mahasiswaService: services.NewMahasiswaService(),
	}
}

// @Summary      Cek status pendataan mahasiswa
// @Description  Cek apakah mahasiswa sudah terdaftar berdasarkan email dan NIK
// @Tags         Mahasiswa
// @Accept       json
// @Produce      json
// @Param        request  body      models.CheckPendataanRequest  true  "Email dan NIK"
// @Success      200      {object}  models.CheckPendataanResponse
// @Failure      400      {object}  map[string]string
// @Router       /mahasiswa/cek-pendataan [post]
func (ctrl *MahasiswaController) CheckPendataan(c *gin.Context) {
	var req models.CheckPendataanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.mahasiswaService.CheckPendataan(req.Email, req.NIK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Registrasi mahasiswa baru
// @Description  Registrasi mahasiswa dengan email, NIK, dan PIN
// @Tags         Mahasiswa
// @Accept       json
// @Produce      json
// @Param        request  body      models.RegisterMahasiswaRequest  true  "Data registrasi"
// @Success      201      {object}  models.RegisterMahasiswaResponse
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Router       /mahasiswa/register [post]
func (ctrl *MahasiswaController) Register(c *gin.Context) {
	var req models.RegisterMahasiswaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.mahasiswaService.Register(&req)
	if err != nil {
		if err.Error() == "Email sudah terdaftar" || err.Error() == "NIK sudah terdaftar" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// @Summary      Verifikasi PIN mahasiswa
// @Description  Verifikasi PIN untuk akses pendataan
// @Tags         Mahasiswa
// @Accept       json
// @Produce      json
// @Param        request  body      models.VerifyPINRequest  true  "Email, NIK, dan PIN"
// @Success      200      {object}  models.VerifyPINResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      423      {object}  map[string]interface{}
// @Router       /mahasiswa/verify-pin [post]
func (ctrl *MahasiswaController) VerifyPIN(c *gin.Context) {
	var req models.VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lockKey := req.Email + ":" + req.NIK
	if locked, remaining := middleware.IsPinLocked(lockKey); locked {
		c.JSON(http.StatusLocked, gin.H{
			"error":  "Terlalu banyak percobaan PIN salah. Akun dikunci sementara.",
			"retry_after_seconds": int(remaining.Seconds()),
		})
		return
	}

	result, err := ctrl.mahasiswaService.VerifyPIN(req.Email, req.NIK, req.PIN)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !result.Valid {
		if locked := middleware.RecordPinFailure(lockKey); locked {
			c.JSON(http.StatusLocked, gin.H{
				"error":  "Terlalu banyak percobaan PIN salah. Akun dikunci selama 15 menit.",
				"retry_after_seconds": 900,
			})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	middleware.ResetPinAttempts(lockKey)
	c.JSON(http.StatusOK, result)
}
