package main

import (
	"crypto/tls"
	"fmt"

	"github.com/zodius/guren"

	"github.com/sirupsen/logrus"
)

func SimpleAuth(request guren.ProxyRequest) bool {
	fmt.Println("user", request.ClaimUser, request.Password)
	fmt.Println("target", request.RequestURI)
	fmt.Println("source", request.SrcAddr)
	return true
}

func main() {
	certificate, err := tls.LoadX509KeyPair("./certificate.pem", "./key.pem")
	if err != nil {
		panic(err)
	}

	config := guren.GurenConfig{
		Protocol:    guren.HTTPS,
		ListenAddr:  ":8080",
		Logger:      logrus.New(),
		AuthFunc:    SimpleAuth,
		Certificate: certificate,
	}

	guren := guren.New(config)
	guren.Start()
}
