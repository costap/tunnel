package envoy

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig()

	if c == nil {
		t.Fail()
	}
}

func TestConfig_AddTCPProxy(t *testing.T) {
	c := NewConfig().AddTCPProxy(TCPProxy{
		Name:         "proxy1",
		ListenerPort: 80,
		ClusterPort:  1080,
	})
	if c == nil {
		t.Fail()
	}
	if c.tcpProxies == nil {
		t.Fail()
	}
	if len(c.tcpProxies) != 1 {
		t.Errorf("tcpProxies len expected 1 and was %v", len(c.tcpProxies))
	}
}

func TestConfig_ToYaml(t *testing.T) {
	c := NewConfig().AddTCPProxy(
		TCPProxy{
			Name:         "proxy1",
			ListenerPort: 80,
			ClusterPort:  10080,
		}).AddTCPProxy(
		TCPProxy{
			Name:         "proxy2",
			ListenerPort: 443,
			ClusterPort:  10443,
		})
	if c == nil {
		t.Fail()
	}
	if c.tcpProxies == nil {
		t.Fail()
	}
	if len(c.tcpProxies) != 2 {
		t.Errorf("tcpProxies len expected 2 and was %v", len(c.tcpProxies))
	}
	yaml := c.ToYaml()
	if yaml != expectedYaml {
		t.Logf("expected:\n%v\n", expectedYaml)
		t.Logf("found:\n%v\n", yaml)
		t.Errorf("yaml not expected")
	}
}

const expectedYaml = `
static_resources:
  listeners:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 80
    filter_chains:
    # Any requests received on this address are sent through this chain of filters
    - filters:
      - name: envoy.filters.network.tcp_proxy
        typed_config:
          "@type": type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy
          stat_prefix: tcp
          cluster: proxy1
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 443
    filter_chains:
    # Any requests received on this address are sent through this chain of filters
    - filters:
      - name: envoy.filters.network.tcp_proxy
        typed_config:
          "@type": type.googleapis.com/envoy.config.filter.network.tcp_proxy.v2.TcpProxy
          stat_prefix: tcp
          cluster: proxy2
  clusters:
  - name: proxy1
    connect_timeout: 1s
    load_assignment:
      cluster_name: proxy1
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 127.0.0.1
                port_value: 10080
  - name: proxy2
    connect_timeout: 1s
    load_assignment:
      cluster_name: proxy2
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 127.0.0.1
                port_value: 10443
admin:
  access_log_path: "/dev/stdout"
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 15000
`
