package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/RicRax/journaLink/model"
)

func RenderHome(db *gorm.DB, c *gin.Context, id uint) {
	sd := model.GetAllDiariesOfUser(db, c, id)
	var sdstring []string

	for _, d := range sd {
		sdstring = append(sdstring, d.Title)
	}

	wn, err := model.GetUsernameFromId(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "could not get username")
		return
	}

	data := struct {
		Username string
		Diaries  []string
	}{
		Username: wn,
		Diaries:  sdstring,
	}

	if sd == nil {
		data.Diaries = append(data.Diaries, "You don't have any diaries!")
	}

	c.HTML(http.StatusOK, "home.html", data)
}

func RenderAddDiary(c *gin.Context) {
	c.HTML(http.StatusOK, "addDiary.html", nil)
}

func RenderDeleteDiary(c *gin.Context) {
	c.HTML(http.StatusOK, "deleteDiary.html", nil)
}

func RenderViewDiary(c *gin.Context) {
	c.HTML(http.StatusOK, "viewDiary.html", nil)
}
