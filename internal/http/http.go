package http

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"strings"
	"sync"
)

type HTTPProxyRequest struct {
	Host         string
	Credential   string
	IsHTTPS      bool
	RawReqHeader bytes.Buffer
}

func ParseRequest(reader *bufio.Reader) (req HTTPProxyRequest, err error) {
	tp := textproto.NewReader(reader)

	// First line: GET /index.html HTTP/1.0
	var requestLine string
	if requestLine, err = tp.ReadLine(); err != nil {
		return
	}

	method, requestURI, _, ok := parseRequestLine(requestLine)
	if !ok {
		err = errors.New("invalid request")
		return
	}

	if method == "CONNECT" {
		req.IsHTTPS = true
		requestURI = "http://" + requestURI
	}

	// get remote host
	uriInfo, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return
	}

	// Subsequent lines: Key: value.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return
	}

	req.Credential = mimeHeader.Get("Proxy-Authorization")

	if uriInfo.Host == "" {
		req.Host = mimeHeader.Get("Host")
	} else {
		if !strings.Contains(uriInfo.Host, ":") {
			req.Host = uriInfo.Host + ":80"
		} else {
			req.Host = uriInfo.Host
		}
	}

	req.RawReqHeader.WriteString(requestLine + "\r\n")
	for k, vs := range mimeHeader {
		for _, v := range vs {
			req.RawReqHeader.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	req.RawReqHeader.WriteString("\r\n")
	return
}

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	method, remain, ok1 := strings.Cut(line, " ")
	requestURI, proto, ok2 := strings.Cut(remain, " ")
	ok = ok1 && ok2
	return
}

func ServeProxy(proxyRequest HTTPProxyRequest, reader *bufio.Reader, conn net.Conn) {
	remoteConn, err := net.Dial("tcp", proxyRequest.Host)
	if err != nil {
		log.Println(err)
		return
	}

	if proxyRequest.IsHTTPS {
		// if https, should sent 200 to client
		_, err = conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		// if not https, should sent the request header to remote
		_, err = proxyRequest.RawReqHeader.WriteTo(remoteConn)
		if err != nil {
			log.Println(err)
			return
		}
	}

	tunnel(reader, conn, remoteConn)

}

func tunnel(brc *bufio.Reader, conn net.Conn, remoteConn net.Conn) {
	defer remoteConn.Close()
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := brc.WriteTo(remoteConn)
		if err != nil {
			log.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(conn, remoteConn)
		if err != nil {
			log.Println(err)
		}
	}()
	wg.Wait()
}
