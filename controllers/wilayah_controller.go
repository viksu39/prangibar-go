package controllers

import (
	"net/http"
	"strconv"

	"prangibar-go/config"
	"prangibar-go/models"

	"github.com/gin-gonic/gin"
)

type WilayahController struct{}

func NewWilayahController() *WilayahController {
	return &WilayahController{}
}

// @Summary      List semua provinsi
// @Tags         Wilayah
// @Success      200  {array}   models.Provinsi
// @Router       /wilayah/provinsi [get]
func (ctrl *WilayahController) GetProvinsi(c *gin.Context) {
	var provinsi []models.Provinsi
	if err := config.DB.Order("kode ASC").Find(&provinsi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, provinsi)
}

// @Summary      List kabupaten/kota berdasarkan provinsi
// @Tags         Wilayah
// @Param        provinsiId  query     int  true  "ID Provinsi"
// @Success      200         {array}   models.KabupatenKota
// @Router       /wilayah/kabupaten-kota [get]
func (ctrl *WilayahController) GetKabupatenKota(c *gin.Context) {
	provinsiID, _ := strconv.Atoi(c.Query("provinsiId"))
	var kabupatenKota []models.KabupatenKota
	if err := config.DB.Where("provinsi_id = ?", provinsiID).Order("kode ASC").Find(&kabupatenKota).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kabupatenKota)
}

// @Summary      List kecamatan berdasarkan kabupaten/kota
// @Tags         Wilayah
// @Param        kabupatenKotaId  query     int  true  "ID Kabupaten/Kota"
// @Success      200             {array}   models.Kecamatan
// @Router       /wilayah/kecamatan [get]
func (ctrl *WilayahController) GetKecamatan(c *gin.Context) {
	kabupatenKotaID, _ := strconv.Atoi(c.Query("kabupatenKotaId"))
	var kecamatan []models.Kecamatan
	if err := config.DB.Where("kabupaten_kota_id = ?", kabupatenKotaID).Order("kode ASC").Find(&kecamatan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kecamatan)
}

// @Summary      List desa berdasarkan kecamatan
// @Tags         Wilayah
// @Param        kecamatanId  query     int  true  "ID Kecamatan"
// @Success      200         {array}   models.Desa
// @Router       /wilayah/desa [get]
func (ctrl *WilayahController) GetDesa(c *gin.Context) {
	kecamatanID, _ := strconv.Atoi(c.Query("kecamatanId"))
	var desa []models.Desa
	if err := config.DB.Where("kecamatan_id = ?", kecamatanID).Order("kode ASC").Find(&desa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, desa)
}
