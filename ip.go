package cfkvadapter

import (
	"io"
	"net/http"
	"time"
)

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
