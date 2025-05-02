// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"

	"project/models"
	user "project/users"
)

var (
	db   *gorm.DB
	rdb  *redis.Client
	rctx = context.Background()
)

func initRedis() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(rctx).Result()
	if err != nil {
		log.Fatal("failed to connect to redis: ", err)
	}
	return rdb
}

func initDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	user := os.Getenv("MYSQL_USER")
	dbname := os.Getenv("MYSQL_DB")
	pass := os.Getenv("MYSQL_PASS")
	hostname := os.Getenv("MYSQL_HOST")

	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, hostname, dbname)

	db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	return db
}

func main() {
	initDB()

	db = initDB()
	db.AutoMigrate(&models.User{})
	// db.AutoMigrate(&models.User{})

	rdb = initRedis()

	initRedis()

	user.Init(db, rdb)

	// user.Init()

	r := gin.Default()

	basepath := r.Group("v1")

	user.RegisteredUserRoute(basepath)

	r.Run(":8080")
}
