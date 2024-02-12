// controllers/authController.go

package controllers

import (
	"io/ioutil"
	"net/http"

	"context"

	"github.com/gin-contrib/sessions"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin" // adjust the import path based on your project structure
	"github.com/shuttlersit/service-desk/backend/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// controllers/googleAuthController.go

// Assuming you have the GoogleCredentials struct defined in your project
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "YOUR_CLIENT_ID",     // Replace with your Client ID
	ClientSecret: "YOUR_CLIENT_SECRET", // Replace with your Client Secret
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func GoogleAuthHandler(c *gin.Context) {

	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// // Google Sign-In route
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleAuthCallback(c *gin.Context) {
	session := sessions.Default(c)
	ctx := context.Background()

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := googleOauthConfig.Client(ctx, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}

	defer response.Body.Close()
	// Process userInfo and create/find user in your DB, then generate a JWT token for the user
	userInfo, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}
	// Within GoogleAuthCallbackHandler in googleAuthController.go
	session.Set("user_info", userInfo) // userInfo being the data retrieved from Google
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// For demonstration, assume we generate a token and user info here
	jwtToken, _ := middleware.GenerateJWT(string(userInfo))

	// Store JWT in session or send back to client directly
	session.Set("jwt_token", jwtToken)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Redirect or respond based on your application needs
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated", "token": jwtToken})
}

func GoogleAuthCallbackHandler(c *gin.Context) {

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()

	userInfo, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}

	// Within GoogleAuthCallbackHandler in googleAuthController.go
	session := sessions.Default(c)
	session.Set("user_info", userInfo) // userInfo being the data retrieved from Google
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Here, you would typically find or create a user in your database.
	// For simplicity, let's generate a JWT token for the user.

	// Assuming the user's email is the subject in your JWT claims
	jwtToken, err := middleware.GenerateJWT(string(userInfo)) // Adjust this according to your JWT generation logic
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken, "user_info": string(userInfo)})
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Implement other auth-related controller functions here

/*
// Google Sign-In callback route
func GoogleCallback(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Fetch user data from Google using the token
	client := googleOauthConfig.Client(c, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()

	// Parse user info from the response (you might want to validate and sanitize this data)
	var googleUser GoogleUserInfo
	if err := json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Check if a user with the same email exists
	var existingUser User
	for _, u := range users {
		if u.Email == googleUser.Email {
			existingUser = u
			break
		}
	}

	if existingUser.ID != 0 {
		// Handle user login or session management here
	} else {
		// Create a new user account with googleUser data
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sign-In with Google Successful", "user": googleUser})
}
*/
