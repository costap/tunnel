package main

import (
	"github.com/costap/tunnel/internal/app/tunnel"
)

func main() {
	t := tunnel.NewTunnel()
	t.CreateHosting()
}