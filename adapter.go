package cfkvadapter

import (
	_ "embed"
	"encoding/json"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

var (
	ID           string // caddy id
	CaddyFile    []byte // /etc/caddy/Caddyfile
	LastUpdateAt int64  // last update at
)

type KVer interface {
	Get(key string) ([]byte, error)
	Delete(key string) error
	Set(key string, value []byte) error
	List(prefix string) ([][]byte, error)
}

type Config struct {
	ApiToken    string `json:"api_token"`
	AccountID   string `json:"account_id"`
	NamespaceID string `json:"namespace_id"`
}

func init() {
	adapter := &Adapter{}
	caddyconfig.RegisterAdapter(adapter.Label(), adapter)

	if err := prepare(); err != nil {
		panic(err)
	}
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

	dispatcher := NewDispatcher(cfg)

	configRaw, warning, err := dispatcher.TouchDefault()
	if err != nil {
		return nil, nil, err
	}
	go dispatcher.Start()
	return configRaw, warning, nil
}
