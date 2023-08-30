package main

import (
	"crypto/tls"
	"fmt"
	"guren"

	"github.com/sirupsen/logrus"
)

func SimpleAuth(credential string) bool {
	username, password, ok := guren.ParseBasicAuth(credential)
	if !ok {
		return false
	}
	fmt.Println(username, password)
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
