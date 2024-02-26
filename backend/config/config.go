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
}

// Dependency represents a table dependency
type Dependency struct {
	Table      string
	Dependency string
}

// os.Getenv()
// readConfig reads configuration from environment variables
func ReadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return nil, err
	}
	config := &Config{
		DBUsername:    os.Getenv("DB_USERNAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		DBSetupStatus: os.Getenv("DB_SETUP_STATUS"),
	}

	// Check if required config fields are set
	if config.DBUsername == "" || config.DBPassword == "" || config.DBHost == "" || config.DBPort == "" || config.DBName == "" {
		return nil, fmt.Errorf("missing required configuration values")
	}

	return config, nil
}

// DBStatusUpdate updates the DB setup status in the environment
func DBStatusUpdate(config *Config) error {
	err := os.Setenv("DB_SETUP_STATUS", config.DBSetupStatus)
	if err != nil {
		log.Fatal(fmt.Sprintln("Error updating DB setup status:", err))
		return err
	}
	return nil
}

type GoogleConfig struct {
	GoogleOAuthConfig *oauth2.Config
	// Other configurations
}

func LoadGoogleConfig() *GoogleConfig {
	return &GoogleConfig{
		GoogleOAuthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}
