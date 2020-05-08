package tunnel

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

type Endpoint struct {
	Host string
	Port int
	User string
}

func NewEndpoint(s string) *Endpoint {
	endpoint := &Endpoint{
		Host: s,
	}

	if parts := strings.Split(endpoint.Host, "@"); len(parts) > 1 {
		endpoint.User = parts[0]
		endpoint.Host = parts[1]
	}

	if parts := strings.Split(endpoint.Host, ":"); len(parts) > 1 {
		endpoint.Host = parts[0]
		endpoint.Port, _ = strconv.Atoi(parts[1])
	}

	return endpoint
}

func (e *Endpoint) String() string {
	if e.Port == 0 {
		return e.Host
	}
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
	Log    *log.Logger
	C      chan *ssh.Client
}

func (t *SSHTunnel) logf(fmt string, args ...interface{}) {
	if t.Log != nil {
		t.Log.Printf(fmt, args)
	}
}

func (t *SSHTunnel) Start() error {
	server, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		t.logf("server dial error: %s", err)
		return err
	}

	t.logf("connected to %s (1 of 2)\n", t.Server.String())

	// Request the remote side to open port 8080 on all interfaces. (0.0.0.0:8080)
	l, err := server.Listen("tcp", t.Remote.String())
	if err != nil {
		log.Fatal("unable to register tcp forward: ", err)
	}
	defer l.Close()

	t.logf("remote is listening to %s (2 of 2)\n", t.Remote.String())

	t.C <- server

	t.logf("starting accept connection from remote...\n")
	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection %w", err)
		}

		t.logf("accepted connection")
		go t.forward(conn)
	}
}

func (t *SSHTunnel) forward(remoteConn net.Conn) {

	localConn, err := net.Dial("tcp", t.Local.String())
	if err != nil {
		t.logf("server dial error: %s", err)
		return
	}

	t.logf("connected to %s\n", t.Local.String())

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			t.logf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func PrivateKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

// NewSSHTunnel creates a new ssh tunnelctl to sshServer and asks the server to start listening on remote interface
// Will forward and connections to local
func NewSSHTunnel(sshServer string, auth ssh.AuthMethod, remote string, local string) *SSHTunnel {
	localEndpoint := NewEndpoint(local)

	server := NewEndpoint(sshServer)
	if server.Port == 0 {
		server.Port = 22
	}

	sshTunnel := &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// Always accept key.
				return nil
			},
			Timeout: 0,
		},
		Local:  localEndpoint,
		Server: server,
		Remote: NewEndpoint(remote),
		C:      make(chan *ssh.Client),
	}

	return sshTunnel
}
