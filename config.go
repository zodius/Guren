package guren

import (
	"crypto/tls"

	"github.com/sirupsen/logrus"
)

type AuthFunc func(ProxyRequest) bool

type GurenConfig struct {
	Protocol    string
	ListenAddr  string
	Certificate tls.Certificate
	Logger      *logrus.Logger
	AuthFunc    AuthFunc
}

func (c *GurenConfig) TLSRequired() bool {
	return c.Protocol == HTTPS || c.Protocol == SOCKS5_TLS
}
