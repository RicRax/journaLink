package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type diaryInfo struct {
	OwnerID int
	DiaryID int
	Title   string
	Body    string
	Shared  []string
}

func main() {
	r := setupRouter()

	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("publicFront/*")

	db, err := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Diary{}, &DiaryAccess{})

	// login route
	r.GET("/login", func(c *gin.Context) {
		handleLogin(c)
	})

	// OAuth2 routes
	r.GET("/oauth/redirect", func(c *gin.Context) {
		handleOAuth(db, c)
	})

	// home route after authentication
	r.GET("/home", func(c *gin.Context) {
		renderHome(db, c)
	})

	// diary endpoints
	r.POST("/diary", func(c *gin.Context) {
		var info diaryInfo
		if err := c.ShouldBindJSON(&info); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			fmt.Println(err)
			return
		}

		if info.DiaryID != 0 {
			updateDiary(db, info, c)
		} else {
			addDiary(db, info, c)
		}
	})

	r.GET("/diary/:id", func(c *gin.Context) {
		getDiary(db, c)
	})

	r.GET("/diary/shared", func(c *gin.Context) {
		getAllSharedDiaries(db, c)
	})

	r.DELETE("/diary/:id", func(c *gin.Context) {
		deleteDiary(db, c)
	})

	// user endpoints
	r.GET("/users/{id}", func(c *gin.Context) {
		getUser(db, c)
	})

	r.POST("/users", func(c *gin.Context) {
		addUser(db, c)
	})

	return r
}
