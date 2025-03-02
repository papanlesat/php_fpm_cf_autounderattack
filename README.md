git clone https://github.com/papanlesat/php_fpm_cf_autounderattack.git
go mod tidy
GOOS=linux GOARCH=amd64 go build -o cf_underattack
