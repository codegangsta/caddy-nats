package caddynats

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(KV{})
}

type KV struct {
	Action     string
	BucketName string
	Key        string

	bucket nats.KeyValue
	app    *App
	logger *zap.Logger
}

func (KV) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.nats_kv",
		New: func() caddy.Module { return new(KV) },
	}
}

func (kv *KV) Provision(ctx caddy.Context) error {
	kv.logger = ctx.Logger(kv)

	natsAppIface, err := ctx.App("nats")
	if err != nil {
		return fmt.Errorf("getting NATS app: %v. Make sure NATS is configured in global options", err)
	}

	kv.app = natsAppIface.(*App)

	return nil
}

func (k KV) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	k.logger.Info("http handler for kv")
	// when its a KV put... execute the next middleware
	// when its a KV get, return immediatly

	// TODO: Figure out what we need to pass in here

	if k.Action == "get" {
		k.logger.Info("performing nats kv get", zap.String("bucket", k.BucketName), zap.String("key", k.Key))
		return k.natsKVGet(k.Key, w)
	}

	//TODO: handle puts
	return next.ServeHTTP(w, r)
}

func (k KV) natsKVGet(key string, w http.ResponseWriter) error {
	bucket, err := k.app.js.KeyValue(k.BucketName)
	if err != nil {
		return err
	}

	entry, err := bucket.Get(k.Key)
	if err != nil {
		return err
	}

	_, err = w.Write(entry.Value())

	return err
}

var (
	_ caddyhttp.MiddlewareHandler = (*KV)(nil)
	_ caddy.Provisioner           = (*KV)(nil)
	_ caddyfile.Unmarshaler       = (*KV)(nil)
)
