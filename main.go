package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

var (
	store             = cookie.NewStore([]byte("secret"))
	sd    sessionData = sessionData{}
)

func main() {
	r := setupRouter()

	initRand()
	sd.init()

	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("front/*")

	db, err := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Diary{}, &DiaryAccess{}, &User{})

	// login routes
	r.GET("/login", func(c *gin.Context) {
		s := sessions.Default(c)

		t := s.Get("token")

		_, ok := sd.authState[t]

		if ok {
			c.Redirect(http.StatusMovedPermanently, "/home")
		}

		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login/authentication", func(c *gin.Context) {
		s := sessions.Default(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "cookie set error")
		}

		u, err := authenticateUser(db, c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "could not authenticate")
			return
		}

		r := randSeq(10)
		sd.authState[r] = u.UID
		s.Set("token", r)
		s.Save()

		c.JSON(http.StatusOK, gin.H{"redirectPath": "/home"})
		return
	})

	// register route
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	// OAuth2 routes
	r.GET("/oauth/redirect", func(c *gin.Context) {
		handleOAuth(db, c)
	})

	// home route after authentication
	r.GET("/home", func(c *gin.Context) {
		s := sessions.Default(c)
		t := s.Get("token")
		id, ok := sd.authState[t]

		if ok {
			renderHome(db, c, id)
		}
	})

	r.GET("/home/addDiary", func(c *gin.Context) {
		renderAddDiary(db, c)
	})

	// diary endpoints
	r.POST("/diary", func(c *gin.Context) {
		s := sessions.Default(c)
		t := s.Get("token")
		id, ok := sd.authState[t]

		if !ok {
			c.JSON(http.StatusInternalServerError, "could not identify token")
		}
		var info diaryInfo
		if err := c.ShouldBindJSON(&info); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			fmt.Println(err)
			return
		}

		info.OwnerID = int(id)

		if info.DiaryID != 0 {
			updateDiary(db, info, c)
		} else {
			addDiary(db, info, c)
		}
	})

	r.GET("/diary/:id", func(c *gin.Context) {
		getDiary(db, c)
	})

	// r.GET("/diary/shared", func(c *gin.Context) {
	// 	getAllSharedDiaries(db, c)
	// })

	r.DELETE("/diary/:id", func(c *gin.Context) {
		deleteDiary(db, c)
	})

	// user endpoints
	r.POST("/user", func(c *gin.Context) {
		addUser(db, c)
	})

	r.GET("/user", func(c *gin.Context) {
		getUser(db, c)
	})

	return r
}
