// main.go
package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"

	mysqldb "project/db"
	"project/models"
	redisdb "project/redis"
	user "project/users"
)

var (
	db   *gorm.DB
	rdb  *redis.Client
	rctx = context.Background()
)

func main() {

	db = mysqldb.InitDB()
	db.AutoMigrate(&models.User{})

	rdb = redisdb.InitRedis()

	user.Init(db, rdb)

	r := gin.Default()
	basepath := r.Group("v1")
	user.RegisteredUserRoute(basepath)

	r.Run(":8080")
}
