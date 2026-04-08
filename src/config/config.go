package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	OauthConfig *oauth2.Config
	ClientURL   string
	// cookie fix
	CookieDomain string
	JWTSecret   string
	StateKey    = "oauth-state"
)

func Load() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "supersecretkey"
	}

	// 3. Configure OAuth
	OauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}

	ClientURL = os.Getenv("CLIENT_URL")
	if ClientURL == "" && os.Getenv("ENV") == "development" {
		ClientURL = "http://localhost:3000"
	}

	// Optional cookie domain. Keep empty for host-only cookies.
	CookieDomain = ".ondc.org"
}
