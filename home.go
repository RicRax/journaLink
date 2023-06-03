package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleHome(db *gorm.DB, c *gin.Context) {
	c.HTML(http.StatusFound, "welcome.html", nil)
}
