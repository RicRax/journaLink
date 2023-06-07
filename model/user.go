package model

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User model for database
type User struct {
	UID      uint   `gorm:"primarykey"`
	Username string `gorm:"uniqueIndex" json:"username"`
	Email    string `                   json:"email"`
	Password string `                   json:"password"`
}

func AddUser(db *gorm.DB, c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		fmt.Println(err)
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create entry"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUser(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func AuthenticateUser(db *gorm.DB, c *gin.Context) (User, error) {
	var u User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, "bad request body")
		return User{}, err
	}

	if err := db.Where("username= ? AND password= ?", u.Username, u.Password).First(&u).Error; err != nil {
		return User{}, err
	}

	return u, nil
}
