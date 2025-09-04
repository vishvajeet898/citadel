package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"

	conf "github.com/Orange-Health/citadel/conf"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPathsToSkipAuth := conf.GetConfig().GetStringSlice("security.url_paths_to_skip_jwt_auth")
		urlPath := c.Request.URL.Path

		if strings.Contains(urlPath, "swagger") {
			c.Next()
			return
		}

		for _, path := range urlPathsToSkipAuth {
			if strings.Contains(path, urlPath) || strings.Contains(urlPath, path) {
				c.Next()
				return
			}
		}

		claims, err := ExtractClaimsFromToken(c)
		if err != nil {
			c.Abort()
			return
		}
		c.Set("JWT_PAYLOAD", claims)
		c.Next()
	}
}

func ExtractClaimsFromToken(c *gin.Context) (jwt.MapClaims, error) {
	authHeader := c.Request.Header.Get("Authorization")

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return nil, errors.New("invalid token")
	}

	tokenString := authHeaderParts[1]

	signingKey := []byte(conf.GetConfig().GetString("security.jwt_secret"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errMessage := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(errMessage)
		}
		return signingKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return nil, errors.New("invalid token")
	}
}
