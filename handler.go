package caddynats

import "github.com/nats-io/nats.go"

type Handler interface {
	Subscribe(conn *nats.Conn) error
	Unsubscribe(conn *nats.Conn) error
}
