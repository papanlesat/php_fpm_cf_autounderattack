package service

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/option"
	"github.com/cloudflare/cloudflare-go/v3/zones"
)

// CloudflareService wraps the Cloudflare v3 client and the zone ID.
type CloudflareService struct {
	client *cloudflare.Client
	zoneID string
}

// NewCloudflareServiceWithZoneGlobal creates a new CloudflareService instance using the Global API Key.
// It reads CF_API_KEY and CF_EMAIL from the environment.
func NewCloudflareServiceWithZoneGlobal(zoneID string) (*CloudflareService, error) {
	apiKey := os.Getenv("CF_API_KEY")
	email := os.Getenv("CF_EMAIL")
	if apiKey == "" || email == "" {
		return nil, fmt.Errorf("CF_API_KEY or CF_EMAIL environment variable is not set")
	}

	// Create the v3 client using API key and email.
	client := cloudflare.NewClient(
		option.WithAPIKey(apiKey),
		option.WithAPIEmail(email),
	)

	return &CloudflareService{
		client: client,
		zoneID: zoneID,
	}, nil
}

// UpdateSecurityLevel updates the zone's security level.
// If mode is "on", it sets the level to "under_attack".
// If mode is "off", it sets the level to "high".

func (c *CloudflareService) UpdateSecurityLevel(mode string) error {
	var body zones.SecurityLevelParam
	switch mode {
	case "on":
		body = zones.SecurityLevelParam{
			ID:    cloudflare.F(zones.SecurityLevelIDSecurityLevel),
			Value: cloudflare.F(zones.SecurityLevelValueUnderAttack),
		}
	case "off":
		body = zones.SecurityLevelParam{
			ID:    cloudflare.F(zones.SecurityLevelIDSecurityLevel),
			Value: cloudflare.F(zones.SecurityLevelValueMedium),
		}
	default:
		return fmt.Errorf("invalid mode: use 'on' or 'off'")
	}

	ctx := context.Background()
	params := zones.SettingEditParams{
		ZoneID: cloudflare.F(c.zoneID),
		Body:   body,
	}

	resp, err := c.client.Zones.Settings.Edit(ctx, "security_level", params)
	if err != nil {
		return fmt.Errorf("error updating zone setting: %v", err)
	}

	fmt.Printf("Cloudflare zone %s updated security level to: %+v\n", c.zoneID, resp)
	return nil
}

// GetSecurityLevel retrieves the current security level from the zone.
func (c *CloudflareService) GetSecurityLevel() (string, error) {
	ctx := context.Background()
	setting, err := c.client.Zones.Settings.Get(ctx, "security_level", zones.SettingGetParams{
		ZoneID: cloudflare.F(c.zoneID),
	})
	if err != nil {
		return "", fmt.Errorf("error retrieving zone settings: %v", err)
	}

	// Try to assert to zones.SecurityLevelValue instead of string.
	if level, ok := setting.Value.(zones.SecurityLevelValue); ok {
		return string(level), nil
	}

	return "", fmt.Errorf("unexpected type for security level: %T", setting.Value)
}
