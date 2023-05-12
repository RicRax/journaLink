package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
		fmt.Println(err)
		return
	}
	if err := db.Create(&entry).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create entry"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, entry)
}

func getDiaryEntry(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	var diary DiaryEntry

	if err := db.First(&diary, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, diary)
}

func updateDiaryEntry(db *gorm.DB, c *gin.Context) {
	var checkEntry DiaryEntry
	var submittedEntry DiaryEntry
	entryID := c.Param("id")

	if err := c.ShouldBindJSON(&submittedEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"wrong body format": err.Error()})
		fmt.Println(err)
		return
	}

	// Find the diary entry with the given ID
	if err := db.Where("ID = ?", entryID).First(&checkEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		fmt.Println(err)
		return
	}

	intEntryID, _ := strconv.Atoi(entryID)
	if submittedEntry.ID != uint(intEntryID) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID of submitted entry does not correspond to path ID"})
		return
	}
	// Update the diary entry with the new data
	db.Save(&submittedEntry)

	c.JSON(http.StatusOK, gin.H{"data": submittedEntry})
}

func deleteDiaryEntry(db *gorm.DB, c *gin.Context) {
	//assert that posted ID is equal to path ID
	id := c.Param("id")

	if err := db.Delete(&DiaryEntry{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diary entry with id " + id + " deleted"})
}
