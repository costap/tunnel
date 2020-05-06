package main

import (
	"github.com/costap/tunnel/internal/pkg/tunnel"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

func main() {

	var running bool = true

	for running {
		tunnel := tunnel.NewSSHTunnel(
			"root@134.122.111.207",
			ssh.Password("pyplwkvp"),
			"0.0.0.0:11443",
			"192.168.0.24:443",
		)

		tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

		go tunnel.Start()
		server := <-tunnel.C
		err := server.Conn.Wait()
		if err != nil {
			tunnel.Log.Printf("server connection closed with error %v", err)
		}
	}
}
