package cfkvadapter

import (
	"context"
	"strconv"

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

func (kv *CFKV) GetInt64(key string) (int64, error) {
	value, err := kv.Get(key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
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

func (kv *CFKV) ListValues(keys []string) (values [][]byte, err error) {
	if len(keys) == 0 {
		return nil, nil
	}
	values = make([][]byte, 0)

	for _, key := range keys {
		value, err := kv.Get(key)
		if err != nil {
			caddy.Log().Error("error getting from kv: ", zap.Error(err))
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func (kv *CFKV) List(prefix string) (keys []string, err error) {
	cursor := ""
	keys = make([]string, 0)

	for {
		resp, err := kv.api.ListWorkersKVKeys(context.Background(),
			cloudflare.AccountIdentifier(kv.cfg.AccountID),
			cloudflare.ListWorkersKVsParams{
				NamespaceID: kv.cfg.NamespaceID,
				Prefix:      prefix,
				Limit:       1000,
				Cursor:      cursor,
			})
		if err != nil {
			caddy.Log().Error("error listing from kv: ", zap.Error(err))
			return nil, err
		}
		for _, result := range resp.Result {
			keys = append(keys, result.Name)
		}
		if resp.HasMorePages() {
			cursor = resp.Cursor
		} else {
			break
		}
	}
	return keys, nil
}
