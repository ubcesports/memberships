package auth

import (
	"os"
	"testing"
)

func TestTrustedFrontendOriginsParsesMultipleOrigins(t *testing.T) {
	t.Setenv("FRONTEND_URL", "http://localhost:3000, https://memberships.ubcesports.ca ,https://lounge.ubcesports.ca/")

	origins := TrustedFrontendOrigins()

	expected := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
		"http://localhost:3000",
		"https://memberships.ubcesports.ca",
		"https://lounge.ubcesports.ca",
	}

	if len(origins) != len(expected) {
		t.Fatalf("expected %d origins, got %d: %v", len(expected), len(origins), origins)
	}

	for index, origin := range expected {
		if origins[index] != origin {
			t.Fatalf("expected origin %d to be %q, got %q", index, origin, origins[index])
		}
	}
}

func TestTrustedFrontendOriginsKeepsLocalDefaultsWithoutEnv(t *testing.T) {
	previousValue, hadValue := os.LookupEnv("FRONTEND_URL")
	if hadValue {
		defer func() {
			_ = os.Setenv("FRONTEND_URL", previousValue)
		}()
	}
	_ = os.Unsetenv("FRONTEND_URL")

	origins := TrustedFrontendOrigins()

	expected := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
	}

	if len(origins) != len(expected) {
		t.Fatalf("expected %d origins, got %d: %v", len(expected), len(origins), origins)
	}

	for index, origin := range expected {
		if origins[index] != origin {
			t.Fatalf("expected origin %d to be %q, got %q", index, origin, origins[index])
		}
	}
}
