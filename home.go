package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getGitHubUsername(accessToken string, c *gin.Context) string {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return "error"
	}

	req.Header.Set("Accept", "application/vnd.github.+json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return "error sending request"
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.String(response.StatusCode, response.Status)
		return "error api sent error"
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return "error"
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return "error"
	}

	return responseData["login"].(string)
}

func renderHome(db *gorm.DB, c *gin.Context, id uint) {
	sd := getAllDiariesOfUser(db, c, id)
	var sdstring []string

	for _, d := range sd {
		sdstring = append(sdstring, d.Title)
	}

	data := struct {
		Diaries []string
	}{
		Diaries: sdstring,
	}

	if sd == nil {
		data.Diaries = append(data.Diaries, "You don't have any diaries!")
	}

	c.HTML(http.StatusOK, "home.html", data)
}

func renderAddDiary(db *gorm.DB, c *gin.Context) {
	c.HTML(http.StatusOK, "addDiary.html", nil)
}
