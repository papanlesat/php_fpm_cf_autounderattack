package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"papanlesat/security/service"

	"github.com/joho/godotenv"
)

func main() {
	// Dapatkan direktori home user
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Gagal mendapatkan direktori home: %v", err)
	}
	// Bangun path ke file .env di ~/.config/cf/.env
	envPath := filepath.Join(homeDir, ".config", "cf", ".env")
	// Muat file .env dari path tersebut
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Tidak dapat memuat file .env dari %s: %v", envPath, err)
	} else {
		log.Printf("Berhasil memuat file .env dari %s", envPath)
	}

	// Threshold penggunaan CPU (dalam persen)
	const cpuThreshold = 40.0

	// Mapping user ke Cloudflare Zone ID.
	// Pastikan environment variable untuk masing-masing user sudah di-set, misalnya:
	// CF_ZONE_LENSAIN=sdsdsd
	// CF_ZONE_KEDAIPE=sdee
	userZoneMapping := map[string]string{
		"lensain+": os.Getenv("CF_ZONE_LENSAIN"),
		"kedaipe+": os.Getenv("CF_ZONE_KEDAIPE"),
	}

	// Dapatkan CF_API_TOKEN dari environment.
	apiToken := os.Getenv("CF_API_TOKEN")
	if apiToken == "" {
		log.Fatal("CF_API_TOKEN belum di-set")
	}

	// Periksa penggunaan CPU php-fpm per user.
	cpuUsage, err := service.CheckFpmCPU()
	if err != nil {
		log.Fatalf("Gagal mengecek penggunaan CPU php-fpm: %v", err)
	}

	// Iterasi setiap user dan periksa apakah penggunaan CPU melebihi threshold.
	for user, usage := range cpuUsage {
		fmt.Printf("User: %s, CPU: %.2f%%\n", user, usage)
		if usage > cpuThreshold {
			zoneID, ok := userZoneMapping[user]
			if !ok || zoneID == "" {
				fmt.Printf("Mapping Cloudflare zone tidak ditemukan untuk user %s\n", user)
				continue
			}

			// Buat instance CloudflareService dengan API token dan zone ID user tersebut.
			cfService, err := service.NewCloudflareServiceWithZone(apiToken, zoneID)
			if err != nil {
				fmt.Printf("Gagal membuat Cloudflare service untuk user %s: %v\n", user, err)
				continue
			}

			// Aktifkan mode "under_attack" untuk zone terkait.
			if err := cfService.UpdateSecurityLevel("on"); err != nil {
				fmt.Printf("Gagal mengaktifkan under_attack mode untuk user %s: %v\n", user, err)
			} else {
				fmt.Printf("Under_attack mode diaktifkan untuk user %s (zone: %s)\n", user, zoneID)
			}
		}
	}
}
