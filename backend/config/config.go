// config/config.go

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config struct to hold database configuration
type Config struct {
	DBUsername    string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	DBSetupStatus string
	DBDSN         string
	DBScript      string
}

// Dependency represents a table dependency
type Dependency struct {
	Table      string
	Dependency string
}

// os.Getenv()
// readConfig reads configuration from environment variables
func ReadConfig() (*Config, error) {
	err := godotenv.Load("/app/.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		DBUsername:    os.Getenv("DB_USERNAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		DBSetupStatus: os.Getenv("DB_SETUP_STATUS"),
		DBDSN:         os.Getenv("DSN"),
		DBScript:      os.Getenv("DB_SCRIPT"),
	}

	// Check if required config fields are set
	if config.DBUsername == "" || config.DBPassword == "" || config.DBHost == "" || config.DBPort == "" || config.DBName == "" {
		log.Fatal(fmt.Errorf("missing required configuration values: %w", err))
		return nil, fmt.Errorf("missing required configuration values")
	}

	return config, nil
}

// DBStatusUpdate updates the DB setup status in the environment
func DBStatusUpdate(config *Config) error {
	if err := os.Setenv("DB_SETUP_STATUS", config.DBSetupStatus); err != nil {
		log.Fatal(fmt.Errorf("error updating DB setup status: %w", err))
		return fmt.Errorf("error updating DB setup status: %w", err)
	}
	return nil
}

type GoogleConfig struct {
	GoogleOAuthConfig *oauth2.Config
	// Other configurations
}

// LoadGoogleConfig loads Google OAuth configuration
func LoadGoogleConfig() (*GoogleConfig, error) {
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if redirectURL == "" || clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing required Google OAuth configuration values")
	}

	return &GoogleConfig{
		GoogleOAuthConfig: &oauth2.Config{
			RedirectURL:  redirectURL,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}, nil
}
