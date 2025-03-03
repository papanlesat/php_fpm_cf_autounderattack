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
	// Load configuration from ~/.config/cf/.env
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Gagal mendapatkan direktori home: %v", err)
	}
	envPath := filepath.Join(homeDir, ".config", "cf", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Tidak dapat memuat file .env dari %s: %v", envPath, err)
	} else {
		log.Printf("Berhasil memuat file .env dari %s", envPath)
	}

	// Load email configuration (jika tersedia)
	emailConfig, err := service.LoadEmailConfig()
	if err != nil {
		log.Printf("Konfigurasi email tidak lengkap: %v. Notifikasi email tidak akan dikirim.", err)
		// Jika tidak tersedia, emailConfig tetap nil
	}

	// Set the CPU usage threshold (in percentage)
	const cpuThreshold = 40.0

	// Mapping of users to their Cloudflare Zone IDs.
	userZoneMapping := map[string]string{
		"lensain+": os.Getenv("CF_ZONE_LENSAIN"),
		"kedaipe+": os.Getenv("CF_ZONE_KEDAIPE"),
		// Tambahkan mapping lain sesuai kebutuhan.
	}

	// Retrieve CPU usage for php-fpm per user.
	cpuUsage, err := service.CheckFpmCPU()
	if err != nil {
		log.Fatalf("Gagal mengecek penggunaan CPU php-fpm: %v", err)
	}

	// Iterate over the CPU usage data for each user.
	for user, usage := range cpuUsage {
		fmt.Printf("User: %s, CPU: %.2f%%\n", user, usage)
		zoneID, ok := userZoneMapping[user]
		if !ok || zoneID == "" {
			fmt.Printf("Mapping Cloudflare zone tidak ditemukan untuk user %s\n", user)
			continue
		}

		// Create a Cloudflare service instance using the global API key method.
		cfService, err := service.NewCloudflareServiceWithZoneGlobal(zoneID)
		if err != nil {
			fmt.Printf("Gagal membuat Cloudflare service untuk user %s: %v\n", user, err)
			continue
		}

		// Retrieve the current security level for the zone.
		currentLevel, err := cfService.GetSecurityLevel()
		if err != nil {
			fmt.Printf("Gagal mengambil security level untuk user %s: %v\n", user, err)
			continue
		}

		// If CPU usage exceeds the threshold, activate under_attack (if not already active).
		if usage > cpuThreshold {
			if currentLevel != "under_attack" {
				if err := cfService.UpdateSecurityLevel("on"); err != nil {
					fmt.Printf("Gagal mengaktifkan under_attack mode untuk user %s: %v\n", user, err)
				} else {
					fmt.Printf("Under_attack mode diaktifkan untuk user %s (zone: %s)\n", user, zoneID)
					// Kirim notifikasi email jika emailConfig tersedia.
					if emailConfig != nil {
						if err := emailConfig.SendNotification(user, zoneID, "on", usage); err != nil {
							fmt.Printf("Gagal mengirim notifikasi email untuk user %s: %v\n", user, err)
						}
					}
				}
			} else {
				fmt.Printf("User %s sudah berada di mode under_attack\n", user)
			}
		} else {
			// If CPU usage is below threshold and mode is still under_attack, disable it (set to "high").
			if currentLevel == "under_attack" {
				if err := cfService.UpdateSecurityLevel("off"); err != nil {
					fmt.Printf("Gagal menonaktifkan under_attack mode untuk user %s: %v\n", user, err)
				} else {
					fmt.Printf("Under_attack mode dinonaktifkan untuk user %s (zone: %s)\n", user, zoneID)
					// Kirim notifikasi email jika emailConfig tersedia.
					if emailConfig != nil {
						if err := emailConfig.SendNotification(user, zoneID, "off", usage); err != nil {
							fmt.Printf("Gagal mengirim notifikasi email untuk user %s: %v\n", user, err)
						}
					}
				}
			} else {
				fmt.Printf("User %s tidak menggunakan mode under_attack, tidak ada perubahan\n", user)
			}
		}
	}
}
