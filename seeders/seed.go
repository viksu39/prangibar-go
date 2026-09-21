package seeders

import (
	"log"

	"prangibar-go/config"
	"prangibar-go/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	var count int64
	config.DB.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
	admin := models.Admin{
		Email:    "admin@prangibar.id",
		Name:     "Administrator",
		Password: string(hashedPassword),
	}
	config.DB.Create(&admin)
	log.Println("Admin seeded: admin@prangibar.id")
}

func SeedWilayah() {
	var count int64
	config.DB.Model(&models.Provinsi{}).Count(&count)
	if count > 0 {
		return
	}

	provinsi := models.Provinsi{Kode: "12", Nama: "Sumatera Utara"}
	config.DB.Create(&provinsi)

	kabupatenKota := []models.KabupatenKota{
		{Kode: "1201", Nama: "Kabupaten Nias", ProvinsiID: provinsi.ID},
		{Kode: "1202", Nama: "Kabupaten Mandailing Natal", ProvinsiID: provinsi.ID},
		{Kode: "1203", Nama: "Kabupaten Tapanuli Selatan", ProvinsiID: provinsi.ID},
		{Kode: "1204", Nama: "Kabupaten Tapanuli Tengah", ProvinsiID: provinsi.ID},
		{Kode: "1205", Nama: "Kabupaten Tapanuli Utara", ProvinsiID: provinsi.ID},
		{Kode: "1206", Nama: "Kabupaten Toba Samosir", ProvinsiID: provinsi.ID},
		{Kode: "1207", Nama: "Kabupaten Labuhan Batu", ProvinsiID: provinsi.ID},
		{Kode: "1208", Nama: "Kabupaten Asahan", ProvinsiID: provinsi.ID},
		{Kode: "1209", Nama: "Kabupaten Simalungun", ProvinsiID: provinsi.ID},
		{Kode: "1210", Nama: "Kabupaten Dairi", ProvinsiID: provinsi.ID},
		{Kode: "1211", Nama: "Kabupaten Karo", ProvinsiID: provinsi.ID},
		{Kode: "1212", Nama: "Kabupaten Deli Serdang", ProvinsiID: provinsi.ID},
		{Kode: "1213", Nama: "Kabupaten Langkat", ProvinsiID: provinsi.ID},
		{Kode: "1214", Nama: "Kabupaten Nias Selatan", ProvinsiID: provinsi.ID},
		{Kode: "1215", Nama: "Kabupaten Humbang Hasundutan", ProvinsiID: provinsi.ID},
		{Kode: "1216", Nama: "Kabupaten Pakpak Bharat", ProvinsiID: provinsi.ID},
		{Kode: "1217", Nama: "Kabupaten Samosir", ProvinsiID: provinsi.ID},
		{Kode: "1218", Nama: "Kabupaten Serdang Bedagai", ProvinsiID: provinsi.ID},
		{Kode: "1219", Nama: "Kabupaten Batu Bara", ProvinsiID: provinsi.ID},
		{Kode: "1220", Nama: "Kabupaten Padang Lawas Utara", ProvinsiID: provinsi.ID},
		{Kode: "1221", Nama: "Kabupaten Padang Lawas", ProvinsiID: provinsi.ID},
		{Kode: "1222", Nama: "Kabupaten Labuhan Batu Selatan", ProvinsiID: provinsi.ID},
		{Kode: "1223", Nama: "Kabupaten Labuhan Batu Utara", ProvinsiID: provinsi.ID},
		{Kode: "1224", Nama: "Kabupaten Nias Utara", ProvinsiID: provinsi.ID},
		{Kode: "1225", Nama: "Kabupaten Nias Barat", ProvinsiID: provinsi.ID},
		{Kode: "1271", Nama: "Kota Medan", ProvinsiID: provinsi.ID},
		{Kode: "1272", Nama: "Kota Pematang Siantar", ProvinsiID: provinsi.ID},
		{Kode: "1273", Nama: "Kota Sibolga", ProvinsiID: provinsi.ID},
		{Kode: "1274", Nama: "Kota Tanjung Balai", ProvinsiID: provinsi.ID},
		{Kode: "1275", Nama: "Kota Binjai", ProvinsiID: provinsi.ID},
		{Kode: "1276", Nama: "Kota Tebing Tinggi", ProvinsiID: provinsi.ID},
		{Kode: "1277", Nama: "Kota Padang Sidempuan", ProvinsiID: provinsi.ID},
		{Kode: "1278", Nama: "Kota Gunungsitoli", ProvinsiID: provinsi.ID},
	}
	config.DB.Create(&kabupatenKota)

	log.Println("Wilayah seeded: 1 provinsi, 33 kabupaten/kota")
}
