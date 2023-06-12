package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"github.com/RicRax/journalink/auth"
	"github.com/RicRax/journalink/model"
)

const clientIDGoogle = "723078198829-4gvefdsns5223ogerdnrvrf9tutrtkri.apps.googleusercontent.com"

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/oauth/redirect",
	ClientID:     clientIDGoogle,
	ClientSecret: auth.CsGoogle,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

var RandomStates []string

type GoogleUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
}

func HandleOAuthGoogle(db *gorm.DB, c *gin.Context) {
	if c.Query("state") != RandomStates[len(RandomStates)-1] {
		fmt.Println("state is not valid")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	token, err := GoogleOauthConfig.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken,
	)
	if err != nil {
		fmt.Println("could not create request")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("could not parse response")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	result := GoogleUserResult{}

	if err := json.Unmarshal(content, &result); err != nil {
		c.JSON(http.StatusInternalServerError, "could not unmarshal response")
		return
	}

	// create token
	s := sessions.Default(c)
	r := auth.RandSeq(10)
	s.Set("token", r)
	s.Save()

	// if necessary create user
	u := model.User{
		Username: result.Name,
	}

	if !model.CheckUserExists(db, result.Name) {
		if err := db.Create(&u).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
			fmt.Println(err)
			return
		}
	} else {
		db.Where("username = ?", u.Username).First(&u)
	}

	// link token to userid in sessionsData
	// auth.SessionsData.AuthState[r] = u.UID
	c.Redirect(http.StatusPermanentRedirect, "/home")
}
