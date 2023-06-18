package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/RicRax/journalink/auth"
	"github.com/RicRax/journalink/model"
	"github.com/RicRax/journalink/routes"
)

var store = cookie.NewStore([]byte(os.Getenv("STORE_KEY")))

func main() {
	r := setupRouter()

	auth.InitRand()

	auth.JwtKey = []byte(os.Getenv("JWT_KEY"))

	auth.CsGoogle = os.Getenv("CS_GOOGLE")

	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("front/*")

	r.Static("/static", "./static")

	db, err := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&model.Diary{}, &model.DiaryAccess{}, &model.User{})

	// login routes
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login/authentication", func(c *gin.Context) {
		auth.Login(db, c)
	})

	r.GET("/refresh", func(c *gin.Context) {
		if err := auth.RefreshToken(c); err != nil {
			c.JSON(http.StatusInternalServerError, "error refreshing token")
		}
		return
	})

	// register route
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	// OAuth2 routes
	r.GET("/oauth", func(c *gin.Context) {
		routes.RandomStates = append(routes.RandomStates, auth.RandSeq(10))
		url := routes.GoogleOauthConfig.AuthCodeURL(routes.RandomStates[len(routes.RandomStates)-1])
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/oauth/redirect", func(c *gin.Context) {
		routes.HandleOAuthGoogle(db, c)
	})

	// home route after authentication
	r.GET("/home", func(c *gin.Context) {
		id, err := auth.AuthenticateJwt(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error while authenticating token")
			return
		}

		routes.RenderHome(db, c, id)
	})

	r.GET("/home/viewDiary", func(c *gin.Context) {
		routes.RenderViewDiary(c)
	})

	r.GET("/home/addDiary", func(c *gin.Context) {
		routes.RenderAddDiary(c)
	})

	r.GET("/home/deleteDiary", func(c *gin.Context) {
		routes.RenderDeleteDiary(c)
	})

	// diary endpoints
	r.POST("/diary", func(c *gin.Context) {
		uid, err := auth.AuthenticateJwt(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error while authenticating token")
			return
		}

		var info model.DiaryInfo
		if err := c.ShouldBindJSON(&info); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			fmt.Println(err)
			return
		}

		info.OwnerID = uid

		if info.DID != 0 {
			model.UpdateDiary(db, info, c)
		} else {
			model.AddDiary(db, info, c)
		}
	})

	r.GET("/diary/:title", func(c *gin.Context) {
		uid, err := auth.AuthenticateJwt(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error while authenticating token")
			return
		}
		model.GetDiary(db, c, uid)
	})

	r.GET("/diary", func(c *gin.Context) {
		id, err := auth.AuthenticateJwt(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error while authenticating token")
			return
		}

		ds, err := model.GetAllDiariesOfUser(db, c, id)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, ds)
	})

	r.DELETE("/diary", func(c *gin.Context) {
		id, err := auth.AuthenticateJwt(c, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error while authenticating token")
			return
		}
		model.DeleteDiary(db, c, id)
	})

	// user endpoints
	r.POST("/user", func(c *gin.Context) {
		model.AddUser(db, c)
	})

	r.GET("/user", func(c *gin.Context) {
		model.GetUser(db, c)
	})

	return r
}
