package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"project/models"

// 	"github.com/joho/godotenv"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var (
// 	db   *gorm.DB
// 	rdb  *redis.Client
// 	rctx = context.Background()
// )

// func initDB() {
// 	var err error
// 	var url string
// 	err = godotenv.Load()
// 	if err != nil {
// 		log.Println("Error is occurred  on .env file please check")
// 	}

// 	user := os.Getenv("MYSQL_USER")
// 	dbname := os.Getenv("MYSQL_DB")
// 	pass := os.Getenv("MYSQL_PASS")
// 	hostname := os.Getenv("MYSQL_HOST")

// 	url = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, hostname, dbname)

// 	db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("failed to connect to database: ", err)
// 	}
// 	db.AutoMigrate(&models.User{})
// }
