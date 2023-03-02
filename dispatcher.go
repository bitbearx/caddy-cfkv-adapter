package cfkvadapter

import (
	"fmt"
	"strings"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"go.uber.org/zap"
)

const (
	CADDY_CADDYFILE  = "CADDY_%s_CADDYFILE"
	CADDY_SRV        = "CADDY_%s_SRV_%s"
	CADDY_SRV_PREFIX = "CADDY_%s_SRV"
	CADDY_UPDATE_AT  = "CADDY_%s_UPDATE_AT"
)

type Dispatcher struct {
	kv KVer

	lastUpdateAt int64
	httpAdapter  caddyfile.Adapter
}

func NewDispatcher(cfg *Config) *Dispatcher {
	var kv KVer = NewKV(cfg)
	var httpAdapter = caddyfile.Adapter{ServerType: httpcaddyfile.ServerType{}}
	return &Dispatcher{
		kv:           kv,
		lastUpdateAt: 0,
		httpAdapter:  httpAdapter,
	}
}

// TouchDefault parse Caddyfile and return
func (d *Dispatcher) TouchDefault() ([]byte, []caddyconfig.Warning, error) {
	err := d.kv.Set(d.keyUpdateAt(), []byte("0"))
	if err != nil {
		return nil, nil, err
	}
	err = d.kv.Set(d.keyCaddyfile(), initCaddyFile)
	if err != nil {
		return nil, nil, err
	}

	return d.httpAdapter.Adapt(initCaddyFile, nil)
}

func (d *Dispatcher) Start() {
	d.loopCheckConfig()
}

func (d *Dispatcher) loopCheckConfig() {
	defer func() {
		if err := recover(); err != nil {
			caddy.Log().Error("loopCheckConfig panic: ", zap.Any("err", err))
		}
	}()

	for {
		select {
		case <-time.After(10 * time.Second):
			cfg, updated, updateAt, err := d.checkConfig()
			if err != nil {
				caddy.Log().Error("error checkConfig: ", zap.Error(err))
				continue
			}
			if updated {
				err = caddy.Load(cfg, false)
				if err != nil {
					caddy.Log().Error("error caddy.Load: ", zap.Error(err))
					continue
				}
				d.lastUpdateAt = updateAt
			}
		}
	}
}

func (d *Dispatcher) checkConfig() (cfg []byte, updated bool, updateAt int64, err error) {
	// check update at
	updateAt, err = d.kv.GetInt64(d.keyUpdateAt())
	if err != nil {
		return
	}
	if updateAt <= d.lastUpdateAt {
		caddy.Log().Info("checkConfig: no update: ", zap.Int64("updateAt", updateAt), zap.Int64("lastUpdateAt", d.lastUpdateAt))
		return
	}
	updated = true

	// check caddyfile
	vhostKeys, err := d.kv.List(d.keySrvPrefix())
	if err != nil {
		caddy.Log().Error("checkConfig: List error: ", zap.Error(err))
		return
	}
	if len(vhostKeys) == 0 {
		caddy.Log().Info("checkConfig: no vhostKeys")
		return
	}

	values, err := d.kv.ListValues(vhostKeys)
	if err != nil {
		caddy.Log().Error("checkConfig: ListValues error: ", zap.Error(err))
		return
	}

	// fetch caddyfile
	caddyfile, err := d.kv.Get(d.keyCaddyfile())
	if err != nil {
		caddy.Log().Error("checkConfig: Get caddyfile error: ", zap.Error(err))
		return
	}
	if len(caddyfile) == 0 {
		return
	}

	// assemble caddyfile
	var sb strings.Builder
	sb.Write(caddyfile)
	sb.Write([]byte("\n"))
	for _, v := range values {
		sb.Write(v)
		sb.Write([]byte("\n"))
	}

	// parse caddyfile
	var warnings []caddyconfig.Warning
	cfg, warnings, err = d.httpAdapter.Adapt(caddyfile, nil)
	if err != nil {
		return
	}
	if len(warnings) > 0 {
		caddy.Log().Warn("checkConfig warnings: ", zap.Any("warnings", warnings))
	}
	return
}

func (d *Dispatcher) keyUpdateAt() string                      { return d.formatKey(CADDY_UPDATE_AT, ID) }
func (d *Dispatcher) keySrvPrefix() string                     { return d.formatKey(CADDY_SRV_PREFIX, ID) }
func (d *Dispatcher) keySRV(srvID string) string               { return d.formatKey(CADDY_SRV, ID, srvID) }
func (d *Dispatcher) keyCaddyfile() string                     { return d.formatKey(CADDY_CADDYFILE, ID) }
func (d *Dispatcher) formatKey(key string, vals ...any) string { return fmt.Sprintf(key, vals...) }
