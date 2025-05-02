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

func getusers(ctx *gin.Context) {
	data, err := rdb.Get(ctx, "users").Result()
	if err == nil {
		var users []models.User
		json.Unmarshal([]byte(data), &users)
		ctx.JSON(http.StatusOK, users)
		return
	}

	var users []models.User
	db.Find(&users)
	rdb.Set(ctx, "users", users, 5*time.Minute)

	ctx.JSON(http.StatusOK, users)
}

func updateuser(ctx *gin.Context) {
	id := ctx.Param("id")

	var existinguser models.User
	if err := db.First(&existinguser, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var input map[string]interface{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
		return
	}

	if err := db.Model(&existinguser).Updates(input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user", "message": err.Error()})
		return
	}

	var users []models.User
	if err := db.Find(&users).Error; err == nil {
		if data, err := json.Marshal(users); err == nil {
			_ = rdb.Set(ctx, "users", data, 5*time.Minute).Err()
		}
	}

	ctx.JSON(http.StatusOK, existinguser)
}

func deleteuser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := db.Delete(&models.User{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not deleted"})
		return
	}

	var users []models.User
	db.Find(&users)
	data, _ := json.Marshal(users)
	rdb.Set(ctx, "users", data, 5*time.Minute)

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func uploadexcel(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Println("file upload error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no file"})
		return
	}

	src, err := file.Open()
	if err != nil {
		log.Println("error opening file:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file not opened"})
		return
	}
	defer src.Close()

	xlFile, err := excelize.OpenReader(src)
	if err != nil {
		log.Println("excel not parsed:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid excel type"})
		return
	}

	rows, err := xlFile.GetRows("uk-500")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read excel rows"})
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

	for i, row := range rows {
		if i == 0 {
			continue
		}
		jobs <- row
	}
	close(jobs)

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no valid data to insert"})
		return
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Println("AutoMigrate error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare database table"})
		return
	}

	batchSize := 100
	if err := db.CreateInBatches(users, batchSize).Error; err != nil {
		log.Println("DB insert error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to DB"})
		return
	}

	jsonData, _ := json.Marshal(users)
	rdb.Set(ctx, "users", jsonData, 5*time.Minute)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data imported successfully",
		"count":   len(users),
	})
}

func RegisteredUserRoute(rg *gin.RouterGroup) {

	userroute := rg.Group("/user")

	userroute.GET("/get", getusers)
	userroute.PUT("/update/:id", updateuser)
	userroute.DELETE("/delete/:id", deleteuser)
	userroute.POST("/upload", uploadexcel)

}
