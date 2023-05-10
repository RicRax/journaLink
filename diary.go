package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type DiaryEntry struct {
	gorm.Model
	ID    string `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
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

func getDiaryEntry(db *gorm.DB, c *gin.Context){ 
		id := c.Param("id")

		var diaries []DiaryEntry

    //err not working 
  if err := db.First(&diaries, id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
        return
    }

		if err := db.Where("ID = ?", id).Find(&diaries).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
			return
		}

		c.JSON(http.StatusOK, diaries)
	}

func updateDiaryEntry(db *gorm.DB, c *gin.Context){ 
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

func deleteDiaryEntry(c *gin.Context) {
	// id := c.Param("id")

	// TODO: Add code to delete the diary entry with the specified ID from the database

	c.JSON(http.StatusOK, gin.H{"message": "Diary entry deleted"})
}
