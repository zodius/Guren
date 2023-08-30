package guren

// supported protocols from libcurl https://curl.haxx.se/libcurl/c/CURLOPT_PROXY.html
const (
	HTTP       = "http"
	HTTPS      = "https"
	SOCKS5     = "socks5"
	SOCKS5_TLS = "socks5-tls"
)
