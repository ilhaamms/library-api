package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ilhaamms/library-api/entity/data"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || len(tokenString) < 7 {
			c.JSON(401, gin.H{
				"message": "authorization required",
			})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]

		claims := &data.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return data.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"message": "invalid token",
			})
			c.Abort()
			return
		}

		var computedHash string
		var bodyBytes []byte

		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			computedHash = computeHash(bodyBytes)
		} else {
			computedHash = computeHash([]byte(tokenString))
		}

		c.Request.Header.Set("X-Request-Hash", computedHash)

		hashHeader := c.GetHeader("X-Request-Hash")

		if hashHeader != computedHash {
			log.Println("hashHeader", hashHeader)
			c.JSON(400, gin.H{
				"message": "data integrity validation failed",
			})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Set("claims", claims)

		c.Next()
	}
}

func computeHash(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
