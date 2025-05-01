package user

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"project/models"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/xuri/excelize/v2"
// )

// // var (
// // 	db   *gorm.DB
// // 	rdb  *redis.Client
// // 	rctx = context.Background()
// // )

// func UploadExcel(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		log.Println("file upload error:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
// 		return
// 	}

// 	src, err := file.Open()
// 	if err != nil {
// 		log.Println("error opening uploaded file:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
// 		return
// 	}
// 	defer src.Close()

// 	xlFile, err := excelize.OpenReader(src)
// 	if err != nil {
// 		log.Println("Excel parse error:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel format"})
// 		return
// 	}

// 	rows, err := xlFile.GetRows("uk-500")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Excel rows"})
// 		return
// 	}

// 	var users []models.User
// 	for i, row := range rows {
// 		if i == 0 || len(row) < 10 {
// 			continue
// 		}
// 		users = append(users, models.User{
// 			FirstName:   row[0],
// 			LastName:    row[1],
// 			CompanyName: row[2],
// 			Address:     row[3],
// 			City:        row[4],
// 			County:      row[5],
// 			Postal:      row[6],
// 			Phone:       row[7],
// 			Email:       row[8],
// 			Web:         row[9],
// 		})
// 	}

// 	if len(users) == 0 {
// 		log.Println("empty slice found")
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "No data found"})
// 		return
// 	}

// 	if err := db.Create(&users).Error; err != nil {
// 		log.Println("DB insert error:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to DB"})
// 		return
// 	}

// 	jsonData, _ := json.Marshal(users)
// 	rdb.Set("users", jsonData, 5*time.Minute)

// 	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully", "count": len(users)})
// }
