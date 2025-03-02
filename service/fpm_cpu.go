package service

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// CheckFpmCPU mengembalikan map dengan key sebagai username dan value sebagai total penggunaan CPU (dalam persentase)
// untuk proses php-fpm.
func CheckFpmCPU() (map[string]float64, error) {
	// Jalankan perintah ps untuk menampilkan kolom USER, %CPU, dan COMMAND.
	cmd := exec.Command("ps", "-eo", "user,pcpu,comm")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Ubah output menjadi string dan pisahkan berdasarkan baris.
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, errors.New("output perintah ps tidak sesuai")
	}

	// Inisialisasi map untuk menyimpan akumulasi penggunaan CPU per user.
	cpuUsage := make(map[string]float64)

	// Lewati baris header dan iterasi setiap baris proses.
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Pisahkan baris berdasarkan spasi.
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		user := fields[0]
		cpuStr := fields[1]
		command := fields[2]

		// Filter hanya untuk proses yang mengandung "php-fpm".
		if !strings.Contains(command, "php-fpm") {
			continue
		}

		// Parsing nilai CPU ke dalam float.
		cpu, err := strconv.ParseFloat(cpuStr, 64)
		if err != nil {
			continue
		}

		// Tambahkan penggunaan CPU ke user terkait.
		cpuUsage[user] += cpu
	}

	return cpuUsage, nil
}
