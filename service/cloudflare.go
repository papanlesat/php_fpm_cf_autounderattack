package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

// CloudflareService wraps the Cloudflare API client and zone configuration.
type CloudflareService struct {
	api    *cloudflare.API
	zoneID string
}

// NewCloudflareServiceWithZoneGlobal creates a new CloudflareService instance using
// the Global API Key and Email (from environment variables CF_API_KEY and CF_EMAIL)
// and the provided zone ID.
func NewCloudflareServiceWithZoneGlobal(zoneID string) (*CloudflareService, error) {
	apiKey := os.Getenv("CF_API_KEY")
	email := os.Getenv("CF_EMAIL")
	if apiKey == "" {
		return nil, errors.New("CF_API_KEY environment variable is not set")
	}
	if email == "" {
		return nil, errors.New("CF_EMAIL environment variable is not set")
	}
	if zoneID == "" {
		return nil, errors.New("zoneID cannot be empty")
	}
	api, err := cloudflare.New(apiKey, email)
	if err != nil {
		return nil, fmt.Errorf("error creating Cloudflare client: %v", err)
	}
	return &CloudflareService{
		api:    api,
		zoneID: zoneID,
	}, nil
}

// UpdateSecurityLevel updates the zone's security level.
// If mode is "on", it sets the level to "under_attack".
// If mode is "off", it sets the level to "high".
func (c *CloudflareService) UpdateSecurityLevel(mode string) error {
	var secLevel string
	switch mode {
	case "on":
		secLevel = "under_attack"
	case "off":
		secLevel = "high"
	default:
		return fmt.Errorf("invalid mode: use 'on' or 'off'")
	}

	ctx := context.Background()
	// Create a zone-level resource container by specifying the Type as "zones" and providing the zone ID.
	rc := &cloudflare.ResourceContainer{
		Type:       "zones",
		Identifier: c.zoneID,
	}

	// Build the update parameters.
	params := cloudflare.UpdateZoneSettingParams{
		Value: secLevel,
	}

	updated, err := c.api.UpdateZoneSetting(ctx, rc, params)
	if err != nil {
		return fmt.Errorf("error updating zone setting: %v", err)
	}

	fmt.Printf("Cloudflare zone %s updated security level to: %v\n", c.zoneID, updated)
	return nil
}

// GetSecurityLevel retrieves the current security level from the zone settings.
func (c *CloudflareService) GetSecurityLevel() (string, error) {
	ctx := context.Background()
	resp, err := c.api.ZoneSettings(ctx, c.zoneID)
	if err != nil {
		return "", fmt.Errorf("error retrieving zone settings: %v", err)
	}

	// Iterate over the Result slice to find the "security_level" setting.
	for _, s := range resp.Result {
		if s.ID == "security_level" {
			level, ok := s.Value.(string)
			if !ok {
				return "", fmt.Errorf("security_level value is not a string")
			}
			return level, nil
		}
	}
	return "", fmt.Errorf("security_level setting not found")
}
