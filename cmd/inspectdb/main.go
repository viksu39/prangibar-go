package main

import (
	"fmt"
	"log"

	"prangibar-go/config"

	"gorm.io/gorm"
)

func main() {
	config.Load()
	db := config.ConnectDB()
	type col struct {
		TableName  string
		ColumnName string
		ColumnType string
		IsNullable string
		ColumnKey  string
	}
	var cols []col
	err := db.Raw(`
SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN ('perusahaan','pendataan','pendataanumkm','mahasiswa','mahasiswas','pendataan_mahasiswa','pendataan_mahasiswas','blok_ii_blocks','blok_iiblocks','blok_iii_anggota_keluarga','blok_iii_anggota_keluargas','admin','apilog')
ORDER BY TABLE_NAME, ORDINAL_POSITION`).Scan(&cols).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range cols {
		fmt.Printf("%-16s %-18s %-20s nullable=%s key=%s\n", c.TableName, c.ColumnName, c.ColumnType, c.IsNullable, c.ColumnKey)
	}

	type idx struct {
		TableName string
		IndexName string
		Seq       int
		ColName   string
		NonUnique int
	}
	var idxs []idx
	err = db.Raw(`
SELECT TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME, NON_UNIQUE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN ('perusahaan','pendataan','pendataanumkm','mahasiswa','mahasiswas','pendataan_mahasiswa','pendataan_mahasiswas','blok_ii_blocks','blok_iiblocks','blok_iii_anggota_keluarga','blok_iii_anggota_keluargas','admin','apilog')
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX`).Scan(&idxs).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- indexes ---")
	for _, i := range idxs {
		fmt.Printf("%-16s %-30s seq=%d col=%s nonunique=%d\n", i.TableName, i.IndexName, i.Seq, i.ColName, i.NonUnique)
	}

	type fk struct {
		TableName        string
		ConstraintName   string
		ColName          string
		RefTableName     string
		RefColName       string
	}
	var fks []fk
	err = db.Raw(`
SELECT TABLE_NAME, CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN ('perusahaan','pendataan','pendataanumkm','mahasiswa','mahasiswas','pendataan_mahasiswa','pendataan_mahasiswas','blok_ii_blocks','blok_iiblocks','blok_iii_anggota_keluarga','blok_iii_anggota_keluargas','admin','apilog')
  AND REFERENCED_TABLE_NAME IS NOT NULL`).Scan(&fks).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- fks ---")
	for _, f := range fks {
		fmt.Printf("%-16s %-40s %s -> %s.%s\n", f.TableName, f.ConstraintName, f.ColName, f.RefTableName, f.RefColName)
	}

	var _ *gorm.DB = db
}
