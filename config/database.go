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
	// Drop legacy FKs that block column/index ALTERs (Error 1553/3780).
	// MySQL requires FK to be dropped before its supporting index can be
	// dropped/recreated. Ignore errors if FK doesn't exist.
	for _, sql := range []string{
		// Legacy FKs that block index rebuild (Error 1553)
		"ALTER TABLE pendataan DROP FOREIGN KEY Pendataan_perusahaanId_fkey",
		"ALTER TABLE pendataan DROP FOREIGN KEY fk_perusahaan_pendataan",
		"ALTER TABLE pendataanumkm DROP FOREIGN KEY PendataanUmkm_perusahaanId_fkey",
		"ALTER TABLE pendataanumkm DROP FOREIGN KEY fk_perusahaan_pendataan_umkm",
		// Pendataan Mahasiswa children: parent.id=bigint unsigned vs child int (Error 3780)
		// GORM table names are pluralized.
		"ALTER TABLE blok_ii_blocks DROP FOREIGN KEY fk_pendataan_mahasiswas_blok_ii",
		"ALTER TABLE blok_iii_anggota_keluargas DROP FOREIGN KEY fk_pendataan_mahasiswas_blok_iii",
		"ALTER TABLE blok_iii_anggota_keluarga DROP FOREIGN KEY fk_pendataan_mahasiswas_blok_iii",
	} {
		db.Exec(sql)
	}
	// Pendataan Mahasiswa children: parent.id was bigint unsigned vs model
	// int32→int. Align PK types BEFORE AutoMigrate (no FK yet on these
	// tables), so GORM creates FKs with compatible types (Error 3780).
	for _, sql := range []string{
		"ALTER TABLE pendataan_mahasiswas MODIFY COLUMN id int NOT NULL AUTO_INCREMENT",
		"ALTER TABLE mahasiswas MODIFY COLUMN id int NOT NULL AUTO_INCREMENT",
		"ALTER TABLE blok_iii_anggota_keluargas MODIFY COLUMN id int NOT NULL AUTO_INCREMENT",
	} {
		db.Exec(sql)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS=1")
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
		&models.BlokIIBlock{},
	)
}
