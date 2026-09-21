package config

import (
	"fmt"
	"log"

	"prangibar-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		App.DBUser, App.DBPassword, App.DBHost, App.DBPort, App.DBDatabase)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Admin{},
		&models.Provinsi{},
		&models.KabupatenKota{},
		&models.Kecamatan{},
		&models.Desa{},
		&models.Perusahaan{},
		&models.Pendataan{},
		&models.PendataanUmkm{},
		&models.ApiLog{},
		&models.Mahasiswa{},
		&models.PendataanMahasiswa{},
		&models.BlokIIIAnggotaKeluarga{},
	)
}
