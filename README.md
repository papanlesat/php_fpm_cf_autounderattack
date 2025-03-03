# Cloudflare Under Attack Mode Auto Switcher

This project is a Go-based tool that monitors the CPU usage of php-fpm processes per user and automatically enables Cloudflare's "I'm Under Attack Mode" when the CPU usage exceeds a defined threshold. Each user can be mapped to a different Cloudflare zone, allowing individualized protection.

## Features

### CPU Monitoring
- Checks CPU usage of php-fpm processes per user.

### Dynamic Cloudflare Mode Switching
- Automatically switches Cloudflare security mode to `under_attack` when a user's CPU consumption goes above a threshold.

### Multi-Zone Support
- Supports mapping different users to different Cloudflare Zone IDs.

## Requirements

### Go
- Go 1.16 or higher is recommended.

### Cloudflare API Token
- A valid Cloudflare API token with permissions to update zone settings.

### System
- Linux/Unix-based OS with the `ps` command available.

## Setup

### Clone the Repository

```bash
git clone https://github.com/<your_username>/<repository_name>.git
cd <repository_name>
```

### Initialize the Go Module and Download Dependencies

```bash
go mod tidy
```

### Configure Environment Variables

Create a `.env` file at `~/.config/cf/.env` with the following content (update with your own values):

```dotenv
# Global API Key credentials
CF_API_KEY=your_global_api_key
CF_EMAIL=your_email@example.com

# Zone IDs per user
CF_ZONE_LENSAIN=sdsdsd
CF_ZONE_KEDAIPE=sdee

# SMTP configuration for email notifications
SMTP_SERVER=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username@example.com
SMTP_PASSWORD=your_smtp_password
ALERT_EMAIL=alert_recipient@example.com

```

You can add more environment variables for additional user-zone mappings as needed.

## Build

To build the application, run the following command:

```bash
go build -o cf_underattack
```

**Tip:** Ensure you are building for the correct target platform. For example, for Linux on `amd64`:

```bash
GOOS=linux GOARCH=amd64 go build -o cf_underattack
```

## Run

After building the binary, execute it:

```bash
./cf_underattack
```

When executed, the application will:
- Load environment variables from `~/.config/cf/.env`.
- Monitor php-fpm CPU usage per user.
- Check if any user's CPU usage exceeds the threshold (default is 40%).
- Activate Cloudflare's `under_attack` mode for the corresponding user's zone if the threshold is surpassed.

## Folder Structure

```go
papanlesat/
├── main.go
└── service/
    ├── cloudflare.go   // Contains Cloudflare API integration logic.
    └── fpm_cpu.go      // Contains the php-fpm CPU monitoring functionality.
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have any improvements or bug fixes.

## License

This project is licensed under the MIT License.
