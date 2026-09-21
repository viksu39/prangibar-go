package controllers

import (
	"net/http"
	"strconv"

	"prangibar-go/models"
	"prangibar-go/services"

	"github.com/gin-gonic/gin"
)

type PendataanController struct {
	pendataanService *services.PendataanService
}

func NewPendataanController() *PendataanController {
	return &PendataanController{
		pendataanService: services.NewPendataanService(),
	}
}

// @Summary      Get info perusahaan via token (publik)
// @Tags         Pendataan
// @Param        token  path      string  true  "Token 64-char hex"
// @Success      200    {object}  models.Perusahaan
// @Failure      404    {object}  map[string]string
// @Router       /pendataan/form/{token} [get]
func (ctrl *PendataanController) GetForm(c *gin.Context) {
	token := c.Param("token")
	perusahaan, err := ctrl.pendataanService.FindByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token tidak valid"})
		return
	}
	c.JSON(http.StatusOK, perusahaan)
}

// @Summary      Submit kuesioner pendataan (publik)
// @Tags         Pendataan
// @Param        token  path      string                 true  "Token 64-char hex"
// @Param        request  body      models.Pendataan     true  "Kuesioner data"
// @Success      201    {object}  models.Pendataan
// @Failure      404    {object}  map[string]string
// @Router       /pendataan/form/{token} [post]
func (ctrl *PendataanController) SubmitForm(c *gin.Context) {
	token := c.Param("token")

	var pendataan models.Pendataan
	if err := c.ShouldBindJSON(&pendataan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.pendataanService.Submit(token, &pendataan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// @Summary      List semua pendataan (admin)
// @Tags         Pendataan
// @Security     BearerAuth
// @Success      200  {array}   models.Pendataan
// @Router       /pendataan [get]
func (ctrl *PendataanController) FindAll(c *gin.Context) {
	pendataan, err := ctrl.pendataanService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pendataan)
}

// @Summary      Detail pendataan (admin)
// @Tags         Pendataan
// @Security     BearerAuth
// @Param        id   path      int  true  "Pendataan ID"
// @Success      200  {object}  models.Pendataan
// @Failure      404  {object}  map[string]string
// @Router       /pendataan/{id} [get]
func (ctrl *PendataanController) FindOne(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	pendataan, err := ctrl.pendataanService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pendataan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, pendataan)
}
