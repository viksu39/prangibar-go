package controllers

import (
	"net/http"
	"strconv"

	"prangibar-go/models"
	"prangibar-go/services"

	"github.com/gin-gonic/gin"
)

type PendataanUmkmController struct {
	pendataanUmkmService *services.PendataanUmkmService
}

func NewPendataanUmkmController() *PendataanUmkmController {
	return &PendataanUmkmController{
		pendataanUmkmService: services.NewPendataanUmkmService(),
	}
}

// @Summary      Get info perusahaan UMKM via token (publik)
// @Tags         Pendataan UMKM
// @Param        token  path      string  true  "Token 64-char hex"
// @Success      200    {object}  models.Perusahaan
// @Failure      404    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Router       /pendataan-umkm/form/{token} [get]
func (ctrl *PendataanUmkmController) GetForm(c *gin.Context) {
	token := c.Param("token")
	perusahaan, err := ctrl.pendataanUmkmService.FindByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perusahaan)
}

// @Summary      Submit kuesioner pendataan UMKM (publik)
// @Tags         Pendataan UMKM
// @Param        token    path      string                    true  "Token 64-char hex"
// @Param        request  body      models.PendataanUmkm     true  "Kuesioner UMKM data"
// @Success      201      {object}  models.PendataanUmkm
// @Failure      404      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Router       /pendataan-umkm/form/{token} [post]
func (ctrl *PendataanUmkmController) SubmitForm(c *gin.Context) {
	token := c.Param("token")

	var pendataan models.PendataanUmkm
	if err := c.ShouldBindJSON(&pendataan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.pendataanUmkmService.Submit(token, &pendataan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// @Summary      List semua pendataan UMKM (admin)
// @Tags         Pendataan UMKM
// @Security     BearerAuth
// @Success      200  {array}   models.PendataanUmkm
// @Router       /pendataan-umkm [get]
func (ctrl *PendataanUmkmController) FindAll(c *gin.Context) {
	pendataan, err := ctrl.pendataanUmkmService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pendataan)
}

// @Summary      Detail pendataan UMKM (admin)
// @Tags         Pendataan UMKM
// @Security     BearerAuth
// @Param        id   path      int  true  "Pendataan UMKM ID"
// @Success      200  {object}  models.PendataanUmkm
// @Failure      404  {object}  map[string]string
// @Router       /pendataan-umkm/{id} [get]
func (ctrl *PendataanUmkmController) FindOne(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	pendataan, err := ctrl.pendataanUmkmService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendataan UMKM tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, pendataan)
}
