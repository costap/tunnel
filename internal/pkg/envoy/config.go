package envoy

import (
	"bytes"
	"text/template"
)

type TCPProxy struct {
	Name                      string
	ListenerPort, ClusterPort int
}
type Config struct {
	tcpProxies []TCPProxy
	adminPort  int
}

func NewConfig() *Config {
	return &Config{tcpProxies: []TCPProxy{}, adminPort: 15000}
}

func (c *Config) TCPProxies() []TCPProxy {
	return c.tcpProxies
}
func (c *Config) AddTCPProxy(proxy TCPProxy) *Config {
	c.tcpProxies = append(c.tcpProxies, proxy)
	return c
}

func (c *Config) AdminPort() int {
	return c.adminPort
}
func (c *Config) SetAdminPort(port int) *Config {
	c.adminPort = port
	return c
}

func (c *Config) ToYaml() string {
	tmpl, err := template.New("envoyConfig").Parse(templateString)
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, c); err != nil {
		panic(err)
	}
	return buff.String()
}

const templateString = `
static_resources:
  listeners:
{{- range .TCPProxies }}
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: {{ .ListenerPort }}
    filter_chains:
    # Any requests received on this address are sent through this chain of filters
    - filters:
      - name: envoy.filters.network.tcp_proxy
        typed_config:
          "@type": type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy
          stat_prefix: tcp
          cluster: {{.Name}}
{{- end }}
  clusters:
{{- range .TCPProxies }}
  - name: {{ .Name }}
    connect_timeout: 1s
    load_assignment:
      cluster_name: {{ .Name }}
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 127.0.0.1
                port_value: {{ .ClusterPort }}
{{- end }}
admin:
  access_log_path: "/dev/stdout"
  address:
    socket_address:
      address: 0.0.0.0
      port_value: {{ .AdminPort }}
`
