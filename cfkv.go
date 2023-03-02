package cfkvadapter

import (
	"context"
	"github.com/caddyserver/caddy/v2"
	cloudflare "github.com/cloudflare/cloudflare-go"
	"go.uber.org/zap"
)

type CFKV struct {
	cfg *Config
	api *cloudflare.API
}

func NewKV(cfg *Config) *CFKV {
	api, err := cloudflare.NewWithAPIToken(cfg.ApiToken)
	if err != nil {
		panic(err)
	}
	return &CFKV{
		cfg: cfg,
		api: api,
	}
}

func (kv *CFKV) Get(key string) ([]byte, error) {
	value, err := kv.api.GetWorkersKV(context.Background(),
		cloudflare.AccountIdentifier(kv.cfg.AccountID),
		cloudflare.GetWorkersKVParams{
			NamespaceID: kv.cfg.NamespaceID,
			Key:         key,
		})
	if err != nil {
		caddy.Log().Error("error getting from kv: ", zap.Error(err))
		return nil, err
	}
	return value, nil
}

func (kv *CFKV) Delete(key string) error {
	resp, err := kv.api.DeleteWorkersKVEntry(context.Background(),
		cloudflare.AccountIdentifier(kv.cfg.AccountID),
		cloudflare.DeleteWorkersKVEntryParams{
			NamespaceID: kv.cfg.NamespaceID,
			Key:         key,
		})
	if err != nil {
		caddy.Log().Error("error deleting from kv: ", zap.Error(err))
		return err
	}
	if resp.Success {
		caddy.Log().Info("successfully deleted from kv: ", zap.String("key", key))
	}
	if len(resp.Errors) > 0 {
		caddy.Log().Error("error deleting from kv: ", zap.String("key", key), zap.Any("errors", resp.Errors))
		return err
	}
	return nil
}

func (kv *CFKV) Set(key string, value []byte) error {
	resp, err := kv.api.WriteWorkersKVEntry(context.Background(),
		cloudflare.AccountIdentifier(kv.cfg.AccountID),
		cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: kv.cfg.NamespaceID,
			Key:         key,
			Value:       value,
		})
	if err != nil {
		caddy.Log().Error("error writing to kv: ", zap.String("key", key), zap.Error(err))
		return err
	}
	if resp.Success {
		caddy.Log().Info("successfully wrote to kv: ", zap.String("key", key))
	}
	if len(resp.Errors) > 0 {
		caddy.Log().Error("error writing to kv: ", zap.String("key", key), zap.Any("errors", resp.Errors))
		return err
	}
	return nil
}

func (kv *CFKV) List(prefix string) ([][]byte, error) {
	cursor := ""
	for {
		resp, err := kv.api.ListWorkersKVKeys(context.Background(),
			cloudflare.AccountIdentifier(kv.cfg.AccountID),
			cloudflare.ListWorkersKVsParams{
				NamespaceID: kv.cfg.NamespaceID,
				Prefix:      prefix,
				Limit:       50,
				Cursor:      cursor,
			})
		if err != nil {
			caddy.Log().Error("error listing from kv: ", zap.Error(err))
			return nil, err
		}
		if resp.HasMorePages() {
			cursor = resp.Cursor
		}
		resp.Messages
	}
}
