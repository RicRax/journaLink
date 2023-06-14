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
	OwnerID uint   `gorm:"column:owner_id"        json:"ownerID"`
	Title   string `                              json:"title"`
	Body    string `                              json:"body"`
}

type DiaryInfo struct {
	OwnerID uint
	DID     uint
	Title   string
	Body    string
	Shared  string
}

// DiaryAccess model for database, determines which users have access to a diary
type DiaryAccess struct {
	FKDiary uint `json:"diaryID"      gorm:"column:fk_diary"`
	FKUser  uint `json:"sharedUserID" gorm:"column:fk_user"`
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

	da := DiaryAccess{
		FKDiary: d.DID,
		FKUser:  d.OwnerID,
	}

	if err := db.Create(&da).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "failed to create access")
		return
	}
	c.JSON(http.StatusOK, d)
}

func GetDiary(db *gorm.DB, c *gin.Context, uid uint) {
	title := c.Param("title")

	var d []Diary

	query2 := `SELECT d.DID, d.title, d.owner_id, d.body 
  FROM diaries d, diary_accesses da 
  WHERE d.DID = da.fk_diary
  AND da.fk_user = ? AND d.title = ?
  `

	if err := db.Raw(query2, uid, title).Scan(&d).Error; err != nil {
		c.JSON(http.StatusBadRequest, "could not get diaries")
		return
	}

	c.JSON(http.StatusOK, d)
}

func GetAllDiariesOfUser(db *gorm.DB, c *gin.Context, uid uint) ([]Diary, error) {
	var d []Diary

	query2 := `SELECT d.DID, d.title, d.owner_id, d.body  
  FROM diaries d, diary_accesses da 
  WHERE d.DID= da.fk_diary
  AND da.fk_user = ?
  `

	if err := db.Raw(query2, uid).Scan(&d).Error; err != nil {
		c.JSON(http.StatusBadRequest, "could not get diaries")
		return nil, err
	}
	return d, nil
}

func UpdateDiary(db *gorm.DB, info DiaryInfo, c *gin.Context) {
	var id uint

	if info.DID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing DiaryID"})
		fmt.Println("missing DiaryID")
		return
	}

	id = info.DID

	var check Diary

	// Find the diary entry with the given ID
	if err := db.Where("DID = ?", id).First(&check).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary entry not found"})
		fmt.Println(err)
		return
	}

	// case inviting another user
	if info.Shared != "" {
		var da DiaryAccess
		da.FKDiary = info.DID
		username, err := GetIdFromUsername(db, info.Shared)
		if err != nil {
			c.JSON(http.StatusBadRequest, "could not find user id")
			return
		}
		da.FKUser = username

		if err := db.Create(&da).Error; err != nil {
			c.JSON(http.StatusBadRequest, "failed to create access")
			return
		}
	} else { // case updating body
		db.Model(&Diary{}).Where("DID = ?", info.DID).Update("Body", info.Body)
	}

	c.JSON(http.StatusOK, gin.H{"udpatedDiary": check})
}

func DeleteDiary(db *gorm.DB, c *gin.Context, id uint) {
	var Title struct {
		Title string
	}

	if err := c.ShouldBindJSON(&Title); err != nil {
		c.JSON(http.StatusBadRequest, "Error bad request")
	}

	q := "DELETE FROM diaries WHERE title = ? AND owner_id = ?"

	if err := db.Exec(q, Title.Title, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, "Diary with this title does not exist")
	}

	// ALSO DELETE FROM DIARY ACCESS

	c.JSON(http.StatusOK, gin.H{"message": "Diary entry with title " + Title.Title + " deleted"})
}
