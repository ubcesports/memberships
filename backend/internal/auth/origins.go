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

	for _, value := range strings.Split(os.Getenv("FRONTEND_URL"), ",") {
		origin := strings.TrimRight(strings.TrimSpace(value), "/")
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}
