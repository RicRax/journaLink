package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	r.POST("/diary", createDiaryEntry)
	r.GET("/diary/:id", getDiaryEntry)
	r.PUT("/diary/:id", updateDiaryEntry)
	r.DELETE("/diary/:id", deleteDiaryEntry)

	r.Run(":8080")
}
