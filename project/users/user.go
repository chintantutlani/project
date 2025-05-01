package user

import (
	"encoding/json"
	"net/http"
	"project/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	// rctx = context.Background()
)

func getUsers(c *gin.Context) {
	data, err := rdb.Get("users").Result()
	if err == nil {
		var users []models.User
		json.Unmarshal([]byte(data), &users)
		c.JSON(http.StatusOK, users)
		return
	}

	var users []models.User
	db.Find(&users)
	rdb.Set("users", users, 5*time.Minute)

	c.JSON(http.StatusOK, users)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user.Address = input.Address
	user.City = input.City
	user.CompanyName = input.CompanyName
	user.County = input.County
	user.Email = input.Email
	user.FirstName = input.FirstName
	user.LastName = input.LastName

	user.Phone = input.Phone
	user.Postal = input.Postal
	user.Web = input.Web
	db.Save(&user)

	var users []models.User
	db.Find(&users)
	data, _ := json.Marshal(users)
	rdb.Set("users", data, 5*time.Minute)

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}

	var users []models.User
	db.Find(&users)
	data, _ := json.Marshal(users)
	rdb.Set("users", data, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func RegisteredUserRoute(rg *gin.RouterGroup) {

	userroute := rg.Group("/user")

	userroute.GET("/get", getUsers)
	userroute.PUT("/uodate/:id", updateUser)
	userroute.DELETE("/delete/:id", deleteUser)
	userroute.POST("/upload", UploadExcel)

}
