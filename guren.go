package guren

import (
	"bufio"
	"crypto/tls"
	"net"

	"github.com/zodius/guren/internal/http"
)

type ProxyRequest struct {
	ClaimUser  string
	Password   string
	SrcAddr    string
	RequestURI string
}

type Guren struct {
	config GurenConfig
}

func New(config GurenConfig) *Guren {
	return &Guren{
		config: config,
	}
}

func (g *Guren) Start() {
	var server net.Listener
	var err error
	if g.config.TLSRequired() {
		server, err = tls.Listen("tcp", g.config.ListenAddr, &tls.Config{
			Certificates: []tls.Certificate{g.config.Certificate},
		})
	} else {
		server, err = net.Listen("tcp", g.config.ListenAddr)
	}
	g.config.Logger.Info("Guren is listening on ", g.config.ListenAddr)

	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go g.handle(conn)
	}
}

func (g *Guren) handle(conn net.Conn) {
	if g.config.Protocol == HTTP || g.config.Protocol == HTTPS {
		g.httpProxy(conn)
	} else {
		g.config.Logger.Debug("Unsupported protocol")
	}
}

func (g *Guren) httpProxy(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	httpProxyRequest, err := http.ParseRequest(reader)
	if err != nil {
		g.config.Logger.Debug(err)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	// auth
	username, password, ok := ParseBasicAuth(httpProxyRequest.Credential)
	if !ok {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}
	proxyRequest := ProxyRequest{
		ClaimUser:  username,
		Password:   password,
		SrcAddr:    conn.RemoteAddr().String(),
		RequestURI: httpProxyRequest.URI,
	}
	if g.config.AuthFunc != nil && !g.config.AuthFunc(proxyRequest) {
		g.config.Logger.Debug("Auth failed")
		conn.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"*\"\r\n\r\n"))
		return
	}

	http.ServeProxy(httpProxyRequest, reader, conn)
}
