package routes

import (
	"prangibar-go/controllers"
	"prangibar-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authController := controllers.NewAuthController()
	adminController := controllers.NewAdminController()
	prelistController := controllers.NewPrelistController()
	pendataanController := controllers.NewPendataanController()
	pendataanUmkmController := controllers.NewPendataanUmkmController()
	pendataanMahasiswaController := controllers.NewPendataanMahasiswaController()
	wilayahController := controllers.NewWilayahController()
	logsController := controllers.NewLogsController()
	mahasiswaController := controllers.NewMahasiswaController()

	public := r.Group("")
	{
		public.POST("/auth/login", authController.Login)
		public.GET("/pendataan/form/:token", pendataanController.GetForm)
		public.POST("/pendataan/form/:token", pendataanController.SubmitForm)
		public.GET("/pendataan-umkm/form/:token", pendataanUmkmController.GetForm)
		public.POST("/pendataan-umkm/form/:token", pendataanUmkmController.SubmitForm)
		public.GET("/wilayah/provinsi", wilayahController.GetProvinsi)
		public.GET("/wilayah/kabupaten-kota", wilayahController.GetKabupatenKota)
		public.GET("/wilayah/kecamatan", wilayahController.GetKecamatan)
		public.GET("/wilayah/desa", wilayahController.GetDesa)
		public.POST("/mahasiswa/cek-pendataan", mahasiswaController.CheckPendataan)
		public.POST("/mahasiswa/register", mahasiswaController.Register)
		public.POST("/mahasiswa/verify-pin", mahasiswaController.VerifyPIN)
		public.GET("/pendataan-mahasiswa", pendataanMahasiswaController.GetPendataan)
		public.GET("/pendataan-mahasiswa/:id", pendataanMahasiswaController.GetByID)
		public.POST("/pendataan-mahasiswa/verify-pin", pendataanMahasiswaController.VerifyPIN)
		public.POST("/pendataan-mahasiswa/draft", pendataanMahasiswaController.SaveDraft)
		public.POST("/pendataan-mahasiswa/submit", pendataanMahasiswaController.Submit)
	}

	admin := r.Group("")
	admin.Use(middleware.JWTAuth())
	{
		admin.GET("/admin", adminController.FindAll)
		admin.GET("/admin/me", adminController.GetMe)
		admin.GET("/admin/:id", adminController.FindOne)
		admin.POST("/admin", adminController.Create)
		admin.PATCH("/admin/:id", adminController.Update)
		admin.DELETE("/admin/:id", adminController.Delete)

		admin.GET("/prelist", prelistController.FindAll)
		admin.GET("/prelist/:id", prelistController.FindOne)
		admin.POST("/prelist", prelistController.Create)
		admin.PUT("/prelist/:id", prelistController.Update)
		admin.DELETE("/prelist/:id", prelistController.Delete)
		admin.POST("/prelist/:id/token", prelistController.GenerateToken)
		admin.POST("/prelist/bulk-token", prelistController.BulkGenerateToken)
		admin.GET("/prelist/import/template", prelistController.DownloadTemplate)
		admin.POST("/prelist/import", prelistController.ImportExcel)

		admin.GET("/pendataan", pendataanController.FindAll)
		admin.GET("/pendataan/:id", pendataanController.FindOne)

		admin.GET("/pendataan-umkm", pendataanUmkmController.FindAll)
		admin.GET("/pendataan-umkm/:id", pendataanUmkmController.FindOne)

		admin.DELETE("/pendataan-mahasiswa/:id", pendataanMahasiswaController.Delete)
		admin.GET("/pendataan-mahasiswa/list", pendataanMahasiswaController.ListWithStatus)

		admin.GET("/logs", logsController.FindAll)
		admin.GET("/logs/summary", logsController.Summary)
	}
}
