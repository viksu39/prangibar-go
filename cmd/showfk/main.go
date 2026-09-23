package main

import (
	"fmt"
	"log"

	"prangibar-go/config"
)

func main() {
	config.Load()
	db := config.ConnectDB()

	var tables []string
	if err := db.Raw(`SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND (TABLE_NAME LIKE '%pendataan%' OR TABLE_NAME LIKE '%blok%' OR TABLE_NAME LIKE '%mahasiswa%') ORDER BY TABLE_NAME`).Scan(&tables).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("tables:", tables)

	// FKs
	rows, err := db.Raw(`
SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE() AND REFERENCED_TABLE_NAME IS NOT NULL
ORDER BY TABLE_NAME`).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("--- all FKs ---")
	for rows.Next() {
		var t, c, col, rt, rc string
		if err := rows.Scan(&t, &c, &col, &rt, &rc); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s | %s | %s -> %s.%s\n", t, c, col, rt, rc)
	}

	// Columns of key tables
	for _, t := range []string{"pendataan_mahasiswas", "blok_iii_anggota_keluargas", "mahasiswas"} {
		crows, err := db.Raw(fmt.Sprintf(`
SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = '%s'
ORDER BY ORDINAL_POSITION`, t)).Rows()
		if err != nil {
			fmt.Printf("ERR %s: %v\n", t, err)
			continue
		}
		fmt.Printf("===== %s columns =====\n", t)
		for crows.Next() {
			var name, ctype, nullable, key string
			if err := crows.Scan(&name, &ctype, &nullable, &key); err != nil {
				fmt.Printf("scan: %v\n", err)
				break
			}
			if key != "" || name == "id" || name == "pendataan_mahasiswa_id" || name == "b1r14" {
				fmt.Printf("  %-30s %-25s key=%s\n", name, ctype, key)
			}
		}
		crows.Close()
	}
}
