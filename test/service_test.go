package test

import (
	"os"
	"testing"

	"papanlesat/security/service"
)

func TestUpdateAndGetSecurityLevel(t *testing.T) {
	// Use environment variable for a dedicated test zone ID.
	zoneID := os.Getenv("CF_TEST_ZONE_ID")
	if zoneID == "" {
		t.Skip("CF_TEST_ZONE_ID is not set; skipping integration test")
	}

	// Create a CloudflareService instance using the global API key method.
	cfService, err := service.NewCloudflareServiceWithZoneGlobal(zoneID)
	if err != nil {
		t.Fatalf("failed to create Cloudflare service: %v", err)
	}

	// Set security level to "high" (which we consider as "off" state) as baseline.
	if err := cfService.UpdateSecurityLevel("off"); err != nil {
		t.Fatalf("failed to set security level to off: %v", err)
	}

	// Retrieve and verify the current security level.
	level, err := cfService.GetSecurityLevel()
	if err != nil {
		t.Fatalf("failed to get security level: %v", err)
	}
	if level != "high" {
		t.Fatalf("expected security level 'high', got: %s", level)
	}

	// Update security level to "under_attack".
	if err := cfService.UpdateSecurityLevel("on"); err != nil {
		t.Fatalf("failed to update security level to on: %v", err)
	}

	// Retrieve and verify the current security level.
	level, err = cfService.GetSecurityLevel()
	if err != nil {
		t.Fatalf("failed to get security level: %v", err)
	}
	if level != "under_attack" {
		t.Fatalf("expected security level 'under_attack', got: %s", level)
	}

	// Revert back to "high" for cleanup.
	if err := cfService.UpdateSecurityLevel("off"); err != nil {
		t.Fatalf("failed to revert security level to off: %v", err)
	}

	level, err = cfService.GetSecurityLevel()
	if err != nil {
		t.Fatalf("failed to get security level after revert: %v", err)
	}
	if level != "high" {
		t.Fatalf("expected security level 'high' after revert, got: %s", level)
	}
}
