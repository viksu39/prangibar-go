package main

import (
	"fmt"
	"log"

	"prangibar-go/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connect root failed:", err)
	}
	name := cfg.DBDatabase
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", name)).Error; err != nil {
		log.Fatal("create db failed:", err)
	}
	fmt.Println("database ready:", name)
}
