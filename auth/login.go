package auth

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/RicRax/journalink/model"
)

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(db *gorm.DB, c *gin.Context) {
	s := sessions.Default(c)

	u, err := AuthenticateUser(db, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "could not authenticate")
		return
	}

	exprirationTime := time.Now().Add(5 * time.Minute)

	claims := Claims{
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exprirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error signing token")
		return
	}

	s.Set("jwtToken", tokenString)

	s.Save()

	c.JSON(http.StatusOK, gin.H{"redirectPath": "/home"})
	return
}

func AuthenticateUser(db *gorm.DB, c *gin.Context) (model.User, error) {
	var u model.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, "bad request body")
		return model.User{}, err
	}

	if err := db.Where("username= ? AND password= ?", u.Username, u.Password).First(&u).Error; err != nil {
		return model.User{}, err
	}

	return u, nil
}
