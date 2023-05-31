package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var userDiaryMap map[int]int = make(map[int]int)

func addUser(db *gorm.DB, c *gin.Context) {
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

func getUser(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func mapUserToDiary(user User, diary Diary) {
	userDiaryMap[int(user.ID)] = int(diary.ID)
}
