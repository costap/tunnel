package tunneld

import (
	"log"
	"os"

	"github.com/costap/tunnel/internal/pkg/tunnel"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	running bool
	started bool
	config  *Config
}

func NewServer(config *Config) *Server {
	return &Server{running: false, started: false, config: config}
}

func (s *Server) Run() {
	s.running = true

	for s.running {
		var auth ssh.AuthMethod
		if s.config.Cert != "" {
			auth = tunnel.PrivateKeyFile(s.config.Cert)
		} else {
			auth = ssh.Password(s.config.Password)
		}
		tunnel := tunnel.NewSSHTunnel(
			s.config.SSHServer,
			auth,
			s.config.RemoteAddr,
			s.config.LocalAddr,
		)

		tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

		go tunnel.Start()
		s.started = true
		sshClient := <-tunnel.C
		err := sshClient.Conn.Wait()
		s.started = false
		if err != nil {
			tunnel.Log.Printf("server connection closed with error %v", err)
		}
	}
}

func (s *Server) Stop() {
	s.running = false

	// TODO: close ssh connection and interrupt main thread.
}

func (s *Server) IsRunning() bool {
	return s.running
}

func (s *Server) IsStarted() bool {
	return s.started
}
