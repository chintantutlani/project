package user

import (
	"encoding/json"
	"log"
	"net/http"
	"project/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	// rctx = context.Background()
)

func Init(database *gorm.DB, redisClient *redis.Client) {
	db = database
	rdb = redisClient
}

func getUsers(c *gin.Context) {
	data, err := rdb.Get(c, "users").Result()
	if err == nil {
		var users []models.User
		json.Unmarshal([]byte(data), &users)
		c.JSON(http.StatusOK, users)
		return
	}

	var users []models.User
	db.Find(&users)
	rdb.Set(c, "users", users, 5*time.Minute)

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
	rdb.Set(c, "users", data, 5*time.Minute)

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
	rdb.Set(c, "users", data, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func UploadExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("file upload error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		log.Println("error opening uploaded file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	xlFile, err := excelize.OpenReader(src)
	if err != nil {
		log.Println("Excel parse error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel format"})
		return
	}

	rows, err := xlFile.GetRows("uk-500")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Excel rows"})
		return
	}

	type result struct {
		user models.User
		ok   bool
	}

	workerCount := 10
	jobs := make(chan []string, len(rows))
	results := make(chan result, len(rows))
	var wg sync.WaitGroup

	// Worker goroutines
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				if len(row) < 10 {
					results <- result{ok: false}
					continue
				}
				user := models.User{
					FirstName:   row[0],
					LastName:    row[1],
					CompanyName: row[2],
					Address:     row[3],
					City:        row[4],
					County:      row[5],
					Postal:      row[6],
					Phone:       row[7],
					Email:       row[8],
					Web:         row[9],
				}
				results <- result{user: user, ok: true}
			}
		}()
	}

	// Skip header and send jobs
	for i, row := range rows {
		if i == 0 {
			continue
		}
		jobs <- row
	}
	close(jobs)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	var users []models.User
	for res := range results {
		if res.ok {
			users = append(users, res.user)
		}
	}

	if len(users) == 0 {
		log.Println("no valid user data found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid data to insert"})
		return
	}

	batchSize := 100
	if err := db.CreateInBatches(users, batchSize).Error; err != nil {
		log.Println("DB insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to DB"})
		return
	}

	// Cache in Redis
	jsonData, _ := json.Marshal(users)
	rdb.Set(c, "users", jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"message": "Data imported successfully",
		"count":   len(users),
	})
}

func RegisteredUserRoute(rg *gin.RouterGroup) {

	userroute := rg.Group("/user")

	userroute.GET("/get", getUsers)
	userroute.PUT("/uodate/:id", updateUser)
	userroute.DELETE("/delete/:id", deleteUser)
	userroute.POST("/upload", UploadExcel)

}
