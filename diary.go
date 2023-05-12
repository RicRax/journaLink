package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type DiaryEntry struct {
	gorm.Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

func addDiaryEntry(db *gorm.DB, c *gin.Context) {
	var entry DiaryEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := db.Create(&entry).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create entry"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func getDiaryEntry(db *gorm.DB, c *gin.Context) {
	id := c.Param("ID")

	var diary []DiaryEntry

	if err := db.First(&diary, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
		return
	}

	c.JSON(http.StatusOK, diary)
}

func updateDiaryEntry(db *gorm.DB, c *gin.Context) {
	var entry DiaryEntry
	entryID := c.Param("id")

	// Find the diary entry with the given ID
	if err := db.Where("id = ?", entryID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		return
	}

	// Bind the updated data from the request body to the DiaryEntry struct
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the diary entry with the new data
	db.Model(&entry).Updates(&entry)

	c.JSON(http.StatusOK, gin.H{"data": entry})
}

func deleteDiaryEntry(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&DiaryEntry{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diary entry with id " + id + " deleted"})
}
