package cfkvadapter

import (
	"io"
	"net/http"
	"os"
)

func prepare() error {
	raw, err := os.ReadFile("/etc/caddy/Caddyfile")
	if err != nil {
		panic(err)
	}

	CaddyFile = raw
	ID = GetPublicIP()
	LastUpdateAt = 0
	return nil
}

func GetPublicIP() string {
	client := &http.Client{
		Timeout: 10,
	}
	resp, err := client.Get("https://api.ipify.org?format=text")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(raw)
}
