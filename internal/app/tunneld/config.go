package tunneld

import (
	"fmt"
	"log"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Mode int

const (
	SSH Mode = iota
	TCP
)

type Config struct {
	Password   string
	Cert       string
	SSHServer  string
	RemoteAddr string
	LocalAddr  string
	AdminPort  int
}

func ConfigInit() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/tunneld/")
	viper.AddConfigPath("$HOME/.tunneld")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v \n", err)
	}

	viper.SetEnvPrefix("tunneld")
	viper.AutomaticEnv()

	viper.RegisterAlias("remote", "remoteAddr")
	viper.RegisterAlias("local", "localAddr")

	flag.StringP("password", "p", "", "password for the ssh server")
	flag.StringP("cert", "c", "", "path to the cert for ssh login")
	flag.String("sshServer", "", "ssh server to open tunnel to")
	flag.String("remoteAddr", "0.0.0.0:8080", "remote interface and port to listen on")
	flag.String("localAddr", "", "local address to connect through the tunnel")
	flag.String("type", "ssh", "type of tunnel, only ssh is currently supported")
	flag.Int("adminPort", 8080, "admin port")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &c
}

func ConfigGet(name string) interface{} {
	return viper.Get(name)
}
