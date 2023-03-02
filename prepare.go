package cfkvadapter

import (
	"io"
	"net/http"
	"os"
	"time"
)

func prepare() error {
	var raw []byte

	_, err := os.Stat("/etc/caddy/Caddyfile")
	if os.IsNotExist(err) {
		raw, err = os.ReadFile("./Caddyfile")
	} else {
		raw, err = os.ReadFile("/etc/caddy/Caddyfile")
	}
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
		Timeout: 10 * time.Second,
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
