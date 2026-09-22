package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"prangibar-go/config"
	"prangibar-go/docs"
	"prangibar-go/middleware"
	"prangibar-go/routes"
	"prangibar-go/seeders"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func parseTrustedProxies() []string {
	env := os.Getenv("TRUSTED_PROXIES")
	if env == "" {
		return []string{"127.0.0.1", "::1"}
	}
	parts := strings.Split(env, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// @title           Backend pendataan perusahaan SE2026-L.UB
// @version         1.0
// @description     Backend pendataan perusahaan berbasis kuesioner SE2026 Kuesioner L (UB & UMKM & Mahasiswa)
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345".
func main() {
	cfg := config.Load()
	db := config.ConnectDB()
	config.AutoMigrate(db)

	seeders.SeedAdmin()
	seeders.SeedWilayah()

	r := gin.Default()
	r.SetTrustedProxies(parseTrustedProxies())

	corsOrigin := cfg.CORSOrigin
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimit())
	r.Use(middleware.APILogger())

	routes.SetupRoutes(r)

	docs.SwaggerInfo.Title = cfg.AppName
	docs.SwaggerInfo.Description = "Backend pendataan perusahaan berbasis kuesioner SE2026 Kuesioner L (UB & UMKM & Mahasiswa)"
	docs.SwaggerInfo.Version = "1.0"

	scheme := "http"
	if os.Getenv("APP_ENV") == "production" {
		scheme = "https"
	}

	hostPort := fmt.Sprintf("localhost:%d", cfg.AppPort)
	if cfg.BaseURL != "" {
		hostPort = cfg.BaseURL
		if strings.HasPrefix(cfg.BaseURL, "https://") {
			hostPort = strings.TrimPrefix(cfg.BaseURL, "https://")
		} else if strings.HasPrefix(cfg.BaseURL, "http://") {
			hostPort = strings.TrimPrefix(cfg.BaseURL, "http://")
		}
	}

	docs.SwaggerInfo.Host = hostPort
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{scheme}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	log.Printf("Server running on http://localhost:%d", cfg.AppPort)
	log.Printf("Swagger docs at http://localhost:%d/swagger/index.html", cfg.AppPort)
	r.Run(addr)
}
