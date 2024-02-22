// middleware/roleMiddleware.go

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/models"
)

// RequireRole checks if the user has the required role to access the route
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assuming you've previously set the claims in the context in your auth middleware
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No claims found"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*models.MyCustomClaims) // Type assert to your custom claims struct
		if !ok || userClaims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden - Insufficient role"})
			c.Abort()
			return
		}

		c.Next() // Proceed to the next middleware/handler
	}
}
