package cfkvadapter

import (
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

const (
	CADDY_CADDYFILE  = "CADDY_%s_CADDYFILE"
	CADDY_SRV        = "CADDY_%s_SRV_%s"
	CADDY_SRV_PREFIX = "CADDY_%s_SRV"
	CADDY_UPDATE_AT  = "CADDY_%s_UPDATE_AT"
)

type Dispatcher struct {
	kv KVer
}

func NewDispatcher(cfg *Config) *Dispatcher {
	var kv KVer = NewKV(cfg)
	return &Dispatcher{
		kv: kv,
	}
}

// TouchDefault parse Caddyfile and return
func (d *Dispatcher) TouchDefault() ([]byte, []caddyconfig.Warning, error) {
	err := d.kv.Set(d.keyUpdateAt(), []byte("0"))
	if err != nil {
		return nil, nil, err
	}
	err = d.kv.Set(d.keyCaddyfile(), CaddyFile)
	if err != nil {
		return nil, nil, err
	}

	return caddyfile.Adapter{ServerType: httpcaddyfile.ServerType{}}.Adapt(CaddyFile, nil)
}

func (d *Dispatcher) Start() {
	go d.loopCheckConfig()
}

func (d *Dispatcher) loopCheckConfig() {
	return
}

func (d *Dispatcher) keyUpdateAt() string                      { return d.formatKey(CADDY_UPDATE_AT, ID) }
func (d *Dispatcher) keySrvPrefix() string                     { return d.formatKey(CADDY_SRV_PREFIX, ID) }
func (d *Dispatcher) keySRV(srvID string) string               { return d.formatKey(CADDY_SRV, ID, srvID) }
func (d *Dispatcher) keyCaddyfile() string                     { return d.formatKey(CADDY_CADDYFILE, ID) }
func (d *Dispatcher) formatKey(key string, vals ...any) string { return fmt.Sprintf(key, vals...) }
