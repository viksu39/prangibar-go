package controllers

import (
	"net/http"

	"prangibar-go/middleware"
	"prangibar-go/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	adminService *services.AdminService
}

func NewAuthController() *AuthController {
	return &AuthController{
		adminService: services.NewAdminService(),
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@prangibar.id"`
	Password string `json:"password" binding:"required,min=6" example:"Admin123!"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// @Summary      Login admin
// @Description  Mengembalikan JWT access_token untuk endpoint yang dilindungi
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials"
// @Success      200      {object}  LoginResponse
// @Failure      401      {object}  map[string]string
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := ctrl.adminService.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	token, err := middleware.GenerateToken(admin.ID, admin.Email, admin.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{AccessToken: token})
}
