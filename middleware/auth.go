package middleware

import (
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

		c.Set("claims", claims)
		c.Next()

	}
}
