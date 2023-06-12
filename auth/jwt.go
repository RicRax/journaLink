package auth

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/RicRax/journalink/model"
)

func AuthenticateJwt(c *gin.Context, db *gorm.DB) (uint, error) {
	s := sessions.Default(c)
	tkn := s.Get("jwtToken")

	claims := &Claims{}
	decodedTkn, err := jwt.ParseWithClaims(
		tkn.(string),
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		},
	)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusBadRequest, "invalid signature")
			return 0, err
		}
		c.JSON(http.StatusBadRequest, "bad request")
		return 0, err
	}

	if !decodedTkn.Valid {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return 0, err
	}

	id, err := model.GetIdFromUsername(db, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "could not get id from username")
		return 0, err
	}

	return id, nil
}

func RefreshToken(c *gin.Context) error {
	s := sessions.Default(c)
	tkn := s.Get("jwtToken")

	claims := &Claims{}
	decodedTkn, err := jwt.ParseWithClaims(
		tkn.(string),
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		},
	)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusBadRequest, "invalid signature")
			return err
		}
		c.JSON(http.StatusBadRequest, "bad request")
		return err
	}

	if !decodedTkn.Valid {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return err
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		c.JSON(http.StatusBadRequest, "not enough time has passed to refresh token")
		return err
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error signing string")
		return err
	}

	s.Set("jwtToken", tokenString)
	return nil
}
