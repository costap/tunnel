package tunneld

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Mode int

const (
	SSH Mode = iota
	TCP
)

type Config interface {
	GetMode() Mode
}

type SSHConfig struct {
	Password   string
	Cert       string
	SSHServer  string
	RemoteAddr string
	LocalAddr  string
}

func (*SSHConfig) GetMode() Mode {
	return SSH
}

func ConfigInit() *SSHConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/tunneld/")
	viper.AddConfigPath("$HOME/.tunneld")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetEnvPrefix("tunneld")
	viper.AutomaticEnv()

	flag.StringP("password", "p", "", "password for the ssh server")
	flag.StringP("cert", "c", "", "path to the cert for ssh login")
	flag.String("sshServer", "", "ssh server to open tunnel to")
	flag.String("remote", "0.0.0.0:8080", "remote interface and port to listen on")
	flag.String("local", "", "local address to connect through the tunnel")
	flag.String("type", "ssh", "type of tunnel, only ssh is currently supported")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	var c SSHConfig
	err = viper.Unmarshal(&c)
	if err != nil {
		t.Fatalf("unable to decode into struct, %v", err)
	}

	return &c
}

func ConfigGet(name string) interface{} {
	return viper.Get(name)
}
