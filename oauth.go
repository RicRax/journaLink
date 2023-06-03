package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const clientID = "03178089f0ff2ea0356d"

type oAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func handleOAuth(db *gorm.DB, c *gin.Context) {
	code := c.Query("code")

	reqURL := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		clientID,
		cs,
		code,
	)

	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Print("could not create HTTP request")
		c.JSON(http.StatusBadRequest, "failed to create http request")
		return
	}

	req.Header.Set("accept", "application/json")

	httpClient := http.Client{}

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Print("could not send HTTP request: %v", err)
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	defer res.Body.Close()

	var AuthResponse oAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&AuthResponse); err != nil {
		fmt.Print("could not parse JSON response")
		c.JSON(http.StatusBadRequest, "error could not parse JSON response")
		return
	}

	// AuthResponse.AccessToken
	c.Redirect(http.StatusFound, "/home?access_token="+AuthResponse.AccessToken)
}
