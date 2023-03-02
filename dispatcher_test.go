package cfkvadapter

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/caddyserver/caddy/v2"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func loadConfig() *Config {
	raw, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

func TestDispatcher_MustDefault(t *testing.T) {
	d := NewDispatcher(loadConfig())
	got, got1, err := d.TouchDefault()

	caddy.Log().Info("got", zap.String("config", string(got)))
	assert.Nil(t, err)
	assert.Equal(t, len(got1), 1)
	assert.NotEmpty(t, got)
}
