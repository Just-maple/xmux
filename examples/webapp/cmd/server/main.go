package main

import (
	"github.com/Just-maple/xmux/examples/webapp/pkg/server"
	"log"
	"time"
)

func main() {
	cfg := server.ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
