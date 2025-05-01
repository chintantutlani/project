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

// 	url = fmt.Sprintf("%s:%s@7@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, hostname, dbname)

// 	db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("failed to connect to database: ", err)
// 	}
// 	db.AutoMigrate(&models.User{})
// }

// func initRedis() {
// 	rdb = redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "",
// 		DB:       0,
// 	})
// 	_, err := rdb.Ping(rctx).Result()
// 	if err != nil {
// 		log.Fatal("failed to connect to redis: ", err)
// 	}
// }

// func uploadExcel(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
// 		return
// 	}
// 	src, err := file.Open()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open uploaded file"})
// 		return
// 	}
// 	defer src.Close()

// 	excel, err := excelize.OpenReader(src)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel file"})
// 		return
// 	}

// 	sheet := excel.GetSheetName(0)
// 	rows, err := excel.GetRows(sheet)
// 	if err != nil || len(rows) < 2 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or empty Excel data"})
// 		return
// 	}

// 	var users []User
// 	for _, row := range rows[1:] {
// 		if len(row) < 3 {
// 			continue
// 		}
// 		age, err := strconv.Atoi(row[2])
// 		if err != nil {
// 			continue
// 		}
// 		users = append(users, User{Name: row[0], Email: row[1], Age: age})
// 	}

// 	go func(users []User) {
// 		db.Create(&users)
// 		data, _ := json.Marshal(users)
// 		rdb.Set(rctx, "users", data, 5*time.Minute)
// 	}(users)

// 	c.JSON(http.StatusOK, gin.H{"message": "Import started"})
// }

// func uploadExcel(c *gin.Context) {
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
// 	rdb.Set(rctx, "users", jsonData, 5*time.Minute)

// 	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully", "count": len(users)})
// }
// func getUsers(c *gin.Context) {
// 	data, err := rdb.Get(rctx, "users").Result()
// 	if err == nil {
// 		var users []models.User
// 		json.Unmarshal([]byte(data), &users)
// 		c.JSON(http.StatusOK, users)
// 		return
// 	}

// 	var users []models.User
// 	db.Find(&users)
// 	rdb.Set(rctx, "users", users, 5*time.Minute)
// 	c.JSON(http.StatusOK, users)
// }

// func updateUser(c *gin.Context) {
// 	id := c.Param("id")
// 	var user models.User
// 	if err := db.First(&user, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	var input models.User
// 	if err := c.BindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	user.Address = input.Address
// 	user.City = input.City
// 	user.CompanyName = input.CompanyName
// 	user.County = input.County
// 	user.Email = input.Email
// 	user.FirstName = input.FirstName
// 	user.LastName = input.LastName

// 	user.Phone = input.Phone
// 	user.Postal = input.Postal
// 	user.Web = input.Web
// 	db.Save(&user)

// 	var users []models.User
// 	db.Find(&users)
// 	data, _ := json.Marshal(users)
// 	rdb.Set(rctx, "users", data, 5*time.Minute)

// 	c.JSON(http.StatusOK, user)
// }

// func deleteUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if err := db.Delete(&models.User{}, id).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
// 		return
// 	}

// 	var users []models.User
// 	db.Find(&users)
// 	data, _ := json.Marshal(users)
// 	rdb.Set(rctx, "users", data, 5*time.Minute)

//		c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
//	}

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

	// db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("failed to connect to database: ", err)
	// }
	// var err error // <-- Don't shadow
	// var err error
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

	// RegisteredUserRoute()

	// r.POST("/upload")
	// r.GET("/users", Getuser)
	// r.PUT("/users/:id", updateUser)
	// r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}
