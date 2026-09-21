package controllers

import (
	"net/http"
	"strconv"

	"prangibar-go/models"
	"prangibar-go/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminController struct {
	adminService *services.AdminService
}

func NewAdminController() *AdminController {
	return &AdminController{
		adminService: services.NewAdminService(),
	}
}

type CreateAdminRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@prangibar.id"`
	Name     string `json:"name" binding:"required" example:"Administrator"`
	Password string `json:"password" binding:"required,min=6" example:"Admin123!"`
}

type UpdateAdminRequest struct {
	Email    string `json:"email" example:"admin@prangibar.id"`
	Name     string `json:"name" example:"Administrator"`
	Password string `json:"password" example:"Admin123!"`
}

// @Summary      List semua admin
// @Tags         Admin
// @Security     BearerAuth
// @Success      200  {array}   models.Admin
// @Router       /admin [get]
func (ctrl *AdminController) FindAll(c *gin.Context) {
	admins, err := ctrl.adminService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, admins)
}

// @Summary      Profil admin yang sedang login
// @Tags         Admin
// @Security     BearerAuth
// @Success      200  {object}  models.Admin
// @Router       /admin/me [get]
func (ctrl *AdminController) GetMe(c *gin.Context) {
	adminID, _ := c.Get("adminID")
	id, _ := adminID.(uint)
	admin, err := ctrl.adminService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, admin)
}

// @Summary      Detail admin
// @Tags         Admin
// @Security     BearerAuth
// @Param        id   path      int  true  "Admin ID"
// @Success      200  {object}  models.Admin
// @Failure      404  {object}  map[string]string
// @Router       /admin/{id} [get]
func (ctrl *AdminController) FindOne(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	admin, err := ctrl.adminService.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, admin)
}

// @Summary      Tambah admin baru
// @Tags         Admin
// @Security     BearerAuth
// @Param        request  body      CreateAdminRequest  true  "Admin data"
// @Success      201      {object}  models.Admin
// @Failure      409      {object}  map[string]string
// @Router       /admin [post]
func (ctrl *AdminController) Create(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	admin := models.Admin{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	created, err := ctrl.adminService.Create(&admin)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah digunakan"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// @Summary      Update admin
// @Tags         Admin
// @Security     BearerAuth
// @Param        id       path      int                  true  "Admin ID"
// @Param        request  body      UpdateAdminRequest  true  "Admin data"
// @Success      200      {object}  models.Admin
// @Failure      404      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Router       /admin/{id} [patch]
func (ctrl *AdminController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hashedPassword string
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		hashedPassword = string(hashed)
	}

	updated, err := ctrl.adminService.Update(uint(id), req.Email, req.Name, hashedPassword)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// @Summary      Hapus admin
// @Tags         Admin
// @Security     BearerAuth
// @Param        id   path      int  true  "Admin ID"
// @Success      200  {object}  models.Admin
// @Failure      404  {object}  map[string]string
// @Router       /admin/{id} [delete]
func (ctrl *AdminController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	deleted, err := ctrl.adminService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, deleted)
}
