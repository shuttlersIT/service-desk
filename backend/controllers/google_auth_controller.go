package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/middleware"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthMainController struct {
	GoogleAuthService *services.GoogleAuthService
}

func NewAGoogleAuthMainController(googleAuthService *services.GoogleAuthService) *GoogleAuthController {
	return &GoogleAuthController{
		GoogleAuthService: googleAuthService,
	}
}

// var c string
var conf *oauth2.Config

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func init() {
	//cid := "946670882701-dcidm9tcfdpcikpbjj8rfsb6uci22o4s.apps.googleusercontent.com"
	//cs := "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"
	cid := os.Getenv("GOOGLE_CLIENT_ID")
	cs := os.Getenv("GOOGLE_CLIENT_SECRET")
	rUrl := os.Getenv("GOOGLE_REDIRECT_URL")

	conf = &oauth2.Config{
		ClientID:     cid,
		ClientSecret: cs,
		RedirectURL:  rUrl,
		//RedirectURL:  "https://intel.shuttlers.africa/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

// IndexHandler handles the login
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// AuthHandler handles authentication of a user and initiates a session.
func (s *GoogleAuthMainController) CustomAuthHandler(c *gin.Context) {
	//Declare shuttlers domain
	//shuttlersDomain := "shuttlers.ng"

	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Login failed. Please try again."})
		return
	}

	userInfo, err := fetchGoogleUserInfo(tok)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process the userInfo to create or update the user in your system
	u, err := s.processUserOAuth(&userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user info"})
		return
	}

	// Generate JWT token for the user
	jwtToken, err := middleware.GenerateJWT(u.Email, "User")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Set JWT token in session or return it in the response
	session.Set("user-id", u.Email)
	session.Set("jwt_token", jwtToken)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User authenticated successfully",
		"token":   jwtToken,
		"user":    u,
	})

}

// LoginHandler handles the login procedure.
func LoginHandler(c *gin.Context) {
	state, err := middleware.GenerateRandomState(32)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Error while generating random data."})
		return
	}
	session := sessions.Default(c)
	session.Set("state", state)
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Error while saving session."})
		return
	}
	link := getLoginURL(state)
	c.HTML(http.StatusOK, "auth.html", gin.H{"link": link})
}

// Logout Handler
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "User Sign out successfully",
	})
}

// processUserOAuth processes the OAuth user data, either registering a new user or updating an existing one.
func (s *GoogleAuthMainController) processUserOAuth(userInfo *models.GoogleUserInfo) (*models.Users, error) {
	fname, lastname := getFirstAndLastName(userInfo.Name)
	//var user models.Users
	user, err := s.GoogleAuthService.AuthDBModel.UserDBModel.GetUserByEmail(userInfo.Email)
	if err != nil {
		// New user registration
		user = &models.Users{
			Email:     userInfo.Email,
			FirstName: fname,
			LastName:  lastname,
			// Populate additional fields as necessary
		}
		if err := s.GoogleAuthService.AuthDBModel.UserDBModel.CreateUser(user); err != nil {
			return nil, err
		}
	} else {
		// Existing user - update details
		user.FirstName = fname
		user.LastName = lastname
		if err := s.GoogleAuthService.AuthDBModel.UserDBModel.UpdateUser(user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
