package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PageData struct {
	Title   string
	Welcome string
	Diaries string
}

func getGitHubUsername(c *gin.Context) string {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return "error"
	}

	req.Header.Set("Accept", "application/vnd.github.+json")
	req.Header.Set("Authorization", "Bearer "+c.Query("access_token")) // not working c.Param
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

func renderHome(db *gorm.DB, c *gin.Context) {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Welcome}}</h1>
	</body>
	</html>
	`
	data := PageData{
		Title: "SapphireHome",

		Welcome: "Welcome, " + getGitHubUsername(c),
	}

	tmpl := template.Must(template.New("home.html").Parse(htmlTemplate))
	err := tmpl.Execute(c.Writer, data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
