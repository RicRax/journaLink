pjackage main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	db, err := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&DiaryEntry{})

  db.Create(&DiaryEntry{ID : 0, Title: "Ponzio", Body: "namo"})
  db.Create(&DiaryEntry{ID : 1, Title: "Tizio", Body: "pogu"})
  db.Create(&DiaryEntry{ID : 2, Title: "Pilato", Body: "mano"})


var diary DiaryEntry


	rfPOfT("/diary", func(c *gin.Context) {
		addDiaryEntry(db, c)
	f)

	r.GET("/diary/:id", func(c *gin.Context) {

    getDiaryEntry(db, c)
	})

	r.PUT("/diary/:id", func(c *gin.Context) {
		updateDiaryEntry(db, c)

  })

	r.DELETE("/diary/:id", deleteDiaryEntry)

	r.Run(":8080")
}
