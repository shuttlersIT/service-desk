package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthorizeRequest is used to authorize a request for a certain end-point group.

func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-token-gen-on-server-side")
		if v == nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "Please login."})
			c.Abort()
		}
		c.Next()
	}
}

func AuthorizeAdminRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminSession := sessions.Default(c)
		v := adminSession.Get("agent-token-gen-on-server-side")
		if v == nil {
			c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{"message": "Please login."})
			c.Abort()
		}
		c.Next()
	}
}

// Implement middleware to check if a request is authenticated
func authenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		userID := uint(claims["userID"].(float64))
		c.Set("userID", userID)
		c.Next()
	}
}
