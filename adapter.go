package cfkvadapter

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
)

//go:embed Caddyfile
var defaultCaddyfile []byte

var (
	ID            string // caddy id
	initCaddyFile []byte // /etc/caddy/Caddyfile
)

type KVer interface {
	Get(key string) ([]byte, error)
	GetInt64(key string) (int64, error)
	Delete(key string) error
	Set(key string, value []byte) error
	List(prefix string) (keys []string, err error)
	ListValues(keys []string) (values [][]byte, err error)
}

type Config struct {
	Caddyfile   string `json:"caddyfile"`
	ID          string `json:"id"`
	ApiToken    string `json:"api_token"`
	AccountID   string `json:"account_id"`
	NamespaceID string `json:"namespace_id"`
}

func init() {
	adapter := &Adapter{}
	caddyconfig.RegisterAdapter(adapter.Label(), adapter)
}

type Adapter struct {
}

func (a *Adapter) Label() string {
	return "cfkv"
}

func (a *Adapter) Adapt(body []byte, options map[string]any) ([]byte, []caddyconfig.Warning, error) {
	cfg := &Config{}
	err := json.Unmarshal(body, cfg)
	if err != nil {
		return nil, nil, err
	}

	if err := a.prepare(cfg); err != nil {
		panic(err)
	}

	dispatcher := NewDispatcher(cfg, caddy.Load)

	configRaw, warning, err := dispatcher.TouchDefault()
	if err != nil {
		return nil, nil, err
	}
	go dispatcher.Start()
	return configRaw, warning, nil
}

func (a *Adapter) prepare(cfg *Config) error {
	var raw []byte

	caddyfilePath := cfg.Caddyfile
	if caddyfilePath == "" {
		raw = defaultCaddyfile
	} else {
		_, err := os.Stat(caddyfilePath)
		if os.IsNotExist(err) {
			panic(err)
		} else {
			raw, err = os.ReadFile(caddyfilePath)
			if err != nil {
				panic(err)
			}
		}
	}

	initCaddyFile = raw
	if cfg.ID == "" {
		ID = GetPublicIP()
	}
	return nil
}
