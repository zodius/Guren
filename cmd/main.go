package main

import (
	"fmt"
	"guren"

	"github.com/sirupsen/logrus"
)

func SimpleAuth(credential string) bool {
	fmt.Println(credential)
	return true
}

func main() {
	config := guren.GurenConfig{
		Protocol:   guren.HTTP,
		ListenAddr: ":8080",
		Logger:     logrus.New(),
		AuthFunc:   SimpleAuth,
	}

	guren := guren.New(config)
	guren.Start()
}
