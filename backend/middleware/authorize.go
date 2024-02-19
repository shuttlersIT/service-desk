// backend/middleware/authorize.go

package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login."})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthorizeAdminRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminSession := sessions.Default(c)
		v := adminSession.Get("agent-token-gen-on-server-side")
		if v == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login."})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthenticateMiddleware is used to check if a request is authenticated
func AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // Replace with your secret key
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

var JwtKey = []byte("your_secret_key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Expect header to be "Bearer <TOKEN>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token is valid; you can optionally add claims to context
		c.Set("userID", claims.Subject)
		c.Next()
	}
}

// CustomClaims includes the authorization claims for the JWT token
type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(email, role string) (string, error) {
	claims := CustomClaims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Service Desk",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	return tokenString, err
}

type Claims struct {
	jwt.StandardClaims
	Email string
}

// GenerateStateOauthCookie generates a random state value and sets it as a cookie.
func GenerateStateOauthCookie(c *gin.Context) string {
	state, _ := GenerateRandomState(32)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true, // Prevents JavaScript access to this cookie to enhance security
		Secure:   true, // Ensure cookie is sent over HTTPS
		Path:     "/",  // Cookie available throughout the domain
	}
	http.SetCookie(c.Writer, &cookie)

	return state
}

// generateRandomState generates a secure random string.
func GenerateRandomState(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Handle error
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateStateOauthCookie checks if the state query parameter matches the state cookie.
func ValidateStateOauthCookie(c *gin.Context) bool {
	stateQuery := c.Query("state")
	stateCookie, err := c.Request.Cookie("oauthstate")
	if err != nil {
		// Handle error: Cookie not found
		return false
	}
	// Validate state value
	return stateQuery == stateCookie.Value
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			return
		}
		tokenString := header[len(BearerSchema):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("your-secret-key"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user ID in context
			userID := claims["userID"].(string)
			c.Set("userID", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err})
			return
		}

		c.Next()
	}
}
