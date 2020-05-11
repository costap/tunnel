package tunneld

import (
	"log"
	"os"
	"time"

	"github.com/costap/tunnel/internal/pkg/tunnel"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	running          bool
	started          bool
	config           *Config
	sshClient        *ssh.Client
	checkStatusSleep time.Duration
}

func NewServer(config *Config) *Server {
	return &Server{running: false, started: false, config: config, checkStatusSleep: 1 * time.Second}
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
		s.sshClient = <-tunnel.C
		err := s.sshClient.Conn.Wait()
		s.started = false
		if err != nil {
			tunnel.Log.Printf("server connection closed with error %v", err)
		}
	}
}

func (s *Server) Stop() {
	s.running = false
	if err := s.sshClient.Close(); err != nil {
		log.Printf("error closing connection %v", err)
	}
}

func (s *Server) IsConnected() bool {
	if _, _, err := s.sshClient.Conn.SendRequest("keepalive@openssh.com", true, nil); err != nil {
		log.Printf("error in keep alice request %v", err)
		return err.Error() != "request failed"
	}
	return true
}

func (s *Server) IsRunning() bool {
	return s.running
}

func (s *Server) IsStarted() bool {
	return s.started
}

func (s *Server) checkStatus() {
	for s.running {
		time.Sleep(s.checkStatusSleep)
		if !s.IsConnected() {
			s.Stop()
			time.Sleep(100 * time.Millisecond)
			s.Run()
		}
	}
}
