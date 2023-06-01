package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Diary struct {
	gorm.Model
	OwnerID int    `json:"ownerID"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

type DiaryAccess struct {
	FK_Diary int    `json:"diaryID"`
	FK_User  string `json:"sharedUserID"`
}

func addDiary(db *gorm.DB, info DiaryInfo, c *gin.Context) {
	var diary Diary

	if diary.Title = info.Title; info.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing Title"})
		fmt.Println("missing Title")
		return
	}

	if diary.OwnerID = info.OwnerID; info.OwnerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing OwnerID"})
		fmt.Println("missing OwnerID")
		return
	}

	if err := db.Create(&diary).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create entry"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, diary)
}

func getDiary(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	var diary Diary

	if err := db.First(&diary, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, diary)
}

func updateDiary(db *gorm.DB, info DiaryInfo, c *gin.Context) {
	var entryID int

	if entryID = info.DiaryID; info.DiaryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing DiaryID"})
		fmt.Println("missing DiaryID")
		return
	}

	var checkDiary Diary

	// Find the diary entry with the given ID
	if err := db.Where("ID = ?", entryID).First(&checkDiary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		fmt.Println(err)
		return
	}

	// Update the diary entry body with the new data
	db.Model(&Diary{}).Where("ID = ?", info.DiaryID).Update("Body", info.Body)

	//Update Access table if necessary using Shared field
	var access DiaryAccess
	access.FK_Diary = info.DiaryID
	if info.Shared != nil {
		for i := 0; i < len(info.Shared); i++ {
			access.FK_User = info.Shared[i]
			if err := db.Create(&access).Error; err != nil {
				c.JSON(http.StatusBadRequest, "failed to create access")
				return
			}
		}
	}

	// Return the updated diary
	db.First(&checkDiary, info.DiaryID)
	var checkAccess DiaryAccess
	db.First(&checkAccess)
	println(checkDiary.ID)
	c.JSON(http.StatusOK, gin.H{"udpatedDiary": checkDiary, "newAccesses": checkAccess})
}

func deleteDiary(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Diary{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diary entry with id " + id + " deleted"})
}
