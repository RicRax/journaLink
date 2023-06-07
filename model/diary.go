package model

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Diary Model for database
type Diary struct {
	DID     uint   `gorm:"primaryKey; column:DID"`
	OwnerID int    `gorm:"column:owner_id"        json:"ownerID"`
	Title   string `                              json:"title"`
	Body    string `                              json:"body"`
}

type DiaryInfo struct {
	OwnerID int
	DiaryID int
	Title   string
	Body    string
	Shared  []string
}

// DiaryAccess model for database, determines which users have access to a diary
type DiaryAccess struct {
	FKDiary int    `json:"diaryID"      gorm:"column:fk_diary"`
	FKUser  string `json:"sharedUserID" gorm:"column:fk_user"`
}

func AddDiary(db *gorm.DB, info DiaryInfo, c *gin.Context) {
	var d Diary

	if d.Title = info.Title; info.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing Title"})
		fmt.Println("missing Title")
		return
	}

	if d.OwnerID = info.OwnerID; info.OwnerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing OwnerID"})
		fmt.Println("missing OwnerID")
		return
	}

	if err := db.Create(&d).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create entry"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, d)
}

func GetDiary(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	var d Diary

	if err := db.First(&d, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get entries"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, d)
}

func GetAllDiariesOfUser(db *gorm.DB, c *gin.Context, id uint) []Diary {
	var d []Diary

	// query := WITH other_diaries AS (
	//  SELECT d.DID, d.title, owns.uid as OwnerID
	//  FROM diary_accesses da
	//  JOIN diaries d ON d.DID = da.fk_diary
	//  JOIN users owns ON d.owner_id = owns.uid
	//  WHERE da.fk_user = ?

	//  ),
	query := `
  WITH my_diaries AS (
  SELECT d.DID, d.title, u.uid as OwnerID
  FROM users u, diaries d
  WHERE u.uid = diaries.owner_id AND u.uid = ?
  )
  SELECT * FROM my_diaries
  `

	db.Raw(query, id).Scan(&d)

	return d
}

func UpdateDiary(db *gorm.DB, info DiaryInfo, c *gin.Context) {
	var id int

	if info.DiaryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing DiaryID"})
		fmt.Println("missing DiaryID")
		return
	}

	id = info.DiaryID

	var check Diary

	// Find the diary entry with the given ID
	if err := db.Where("DID = ?", id).First(&check).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		fmt.Println(err)
		return
	}

	// Update the diary entry body with the new data
	db.Model(&Diary{}).Where("DID = ?", info.DiaryID).Update("Body", info.Body)

	// Update Access table if necessary using Shared field
	var da DiaryAccess
	da.FKDiary = info.DiaryID
	if info.Shared != nil {
		for i := 0; i < len(info.Shared); i++ {
			da.FKUser = info.Shared[i]
			if err := db.Create(&da).Error; err != nil {
				c.JSON(http.StatusBadRequest, "failed to create access")
				return
			}
		}
	}

	// Return the updated diary
	db.First(&check, info.DiaryID)
	var checkA DiaryAccess
	db.First(&checkA)
	c.JSON(http.StatusOK, gin.H{"udpatedDiary": check, "newAccesses": checkA})
}

func DeleteDiary(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Diary{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diary entry with id " + id + " deleted"})
}
