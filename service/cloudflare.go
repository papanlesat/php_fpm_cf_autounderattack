package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

// CloudflareService membungkus API Cloudflare dan konfigurasi zone.
type CloudflareService struct {
	api    *cloudflare.API
	zoneID string
}

// NewCloudflareServiceFromEnv membuat instance menggunakan variabel lingkungan (CF_API_TOKEN dan CF_ZONE_ID).
func NewCloudflareServiceFromEnv() (*CloudflareService, error) {
	apiToken := os.Getenv("CF_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("CF_API_TOKEN environment variable is not set")
	}

	zoneID := os.Getenv("CF_ZONE_ID")
	if zoneID == "" {
		return nil, fmt.Errorf("CF_ZONE_ID environment variable is not set")
	}

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Cloudflare client: %v", err)
	}

	return &CloudflareService{
		api:    api,
		zoneID: zoneID,
	}, nil
}

// NewCloudflareServiceWithZone membuat instance CloudflareService dengan zone ID yang ditentukan.
func NewCloudflareServiceWithZone(apiToken, zoneID string) (*CloudflareService, error) {
	if apiToken == "" {
		return nil, errors.New("apiToken tidak boleh kosong")
	}
	if zoneID == "" {
		return nil, errors.New("zoneID tidak boleh kosong")
	}
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Cloudflare client: %v", err)
	}
	return &CloudflareService{
		api:    api,
		zoneID: zoneID,
	}, nil
}

// UpdateSecurityLevel mengubah security level zone menjadi "under_attack" jika mode "on" atau mengembalikannya ke "high" jika mode "off".
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
	// Membuat resource container dengan zone ID.
	rc := &cloudflare.ResourceContainer{Identifier: c.zoneID}
	// Parameter update hanya memerlukan nilai baru.
	params := cloudflare.UpdateZoneSettingParams{
		Value: secLevel,
	}

	updated, err := c.api.UpdateZoneSetting(ctx, rc, params)
	if err != nil {
		return fmt.Errorf("error updating zone setting: %v", err)
	}

	fmt.Printf("Cloudflare zone %s updated security level to: %s\n", c.zoneID, updated.Value)
	return nil
}
