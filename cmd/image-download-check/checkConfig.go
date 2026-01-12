package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// CheckConfig stores configuration for the image download check.
type CheckConfig struct {
	// FullImageURL is the registry/image:tag reference to pull.
	FullImageURL string
	// TimeoutLimit is the maximum duration for the pull.
	TimeoutLimit time.Duration
	// LoginRequired indicates whether registry authentication is needed.
	LoginRequired bool
	// RegistryUsername is the optional registry username.
	RegistryUsername string
	// RegistryPassword is the optional registry password.
	RegistryPassword string
}

// parseConfig loads the environment variables into a configuration struct.
func parseConfig() (*CheckConfig, error) {
	// Read required image reference.
	fullImageURL := os.Getenv("FULL_IMAGE_URL")
	if len(fullImageURL) == 0 {
		return nil, fmt.Errorf("no FULL_IMAGE_URL string provided in YAML")
	}

	// Read the timeout limit.
	timeoutLimitRaw := os.Getenv("TIMEOUT_LIMIT")
	if len(timeoutLimitRaw) == 0 {
		return nil, fmt.Errorf("no TIMEOUT_LIMIT string provided in YAML")
	}

	// Parse the timeout duration.
	timeoutLimit, err := time.ParseDuration(timeoutLimitRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TIMEOUT_LIMIT: %w", err)
	}

	// Determine if registry login is required.
	loginRequired := false
	if strings.ToLower(os.Getenv("LOGIN_REQUIRED")) == "true" {
		loginRequired = true
	}

	// Read optional registry credentials.
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")

	// Assemble configuration.
	cfg := &CheckConfig{}
	cfg.FullImageURL = fullImageURL
	cfg.TimeoutLimit = timeoutLimit
	cfg.LoginRequired = loginRequired
	cfg.RegistryUsername = username
	cfg.RegistryPassword = password

	return cfg, nil
}
