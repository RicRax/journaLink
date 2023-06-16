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
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		fmt.Println(err)
		return
	}

	if u.Username == "" || u.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing username or password"})
		return
	}

	if err := db.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create entry"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, u)
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

func GetUsernameFromId(db *gorm.DB, id uint) (string, error) {
	var u User

	if err := db.Where("UID = ?", id).First(&u).Error; err != nil {
		return "", err
	}

	return u.Username, nil
}

func GetIdFromUsername(db *gorm.DB, username string) (uint, error) {
	var u User

	if err := db.Where("username = ?", username).First(&u).Error; err != nil {
		return 0, err
	}

	return u.UID, nil
}

func CheckUserExists(db *gorm.DB, username string) bool {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err == nil {
		// User exists
		return true
	} else if err == gorm.ErrRecordNotFound {
		// User does not exist
		return false
	} else {
		// Error occurred while querying the database
		// Handle the error
		return false
	}
}
