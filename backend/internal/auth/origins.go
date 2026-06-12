package auth

import (
	"os"
	"strings"
)

func TrustedFrontendOrigins() []string {
	origins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
	}

	if frontendURL := strings.TrimRight(os.Getenv("FRONTEND_URL"), "/"); frontendURL != "" {
		origins = append(origins, frontendURL)
	}

	return origins
}
