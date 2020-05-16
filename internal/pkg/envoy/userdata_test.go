package envoy

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestConfig_CloudConfigYaml(t *testing.T) {
	c := NewConfig().AddTCPProxy(TCPProxy{
		Name:         "proxy1",
		ListenerPort: 80,
		ClusterPort:  1080,
	})
	yaml := c.CloudConfigYaml()
	if yaml != cloudConfigExpected {
		if err := ioutil.WriteFile("./TestConfig_CloudConfigYaml_expected.txt", []byte(cloudConfigExpected), 0644); err != nil {
			fmt.Printf("error writing expected to file: %v", err)
		}
		if err := ioutil.WriteFile("./TestConfig_CloudConfigYaml_found.txt", []byte(yaml), 0644); err != nil {
			fmt.Printf("error writing actual to file: %v", err)
		}
		t.Errorf("expected:[\n%v\n]\nfound:[\n%v\n]", cloudConfigExpected, yaml)
	}
}

const cloudConfigExpected = `#cloud-config
write_files:
- content: |
    *         hard    nofile      500000
    *         soft    nofile      500000
    root      hard    nofile      500000
    root      soft    nofile      500000
  path: /etc/security/limits.conf
- content: "PermitRootLogin yes\nPasswordAuthentication no\nChallengeResponseAuthentication
    no\nUsePAM yes\nX11Forwarding yes\nPrintMotd no\nClientAliveInterval 120\nClientAliveCountMax
    720\nAcceptEnv LANG LC_*\nSubsystem\tsftp\t/usr/lib/openssh/sftp-server"
  path: /etc/ssh/sshd_config
- content: |2

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
                    port_value: 1080
    admin:
      access_log_path: "/dev/stdout"
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 15000
  path: /etc/envoy/config.yaml
runcmd:
- - sh
  - -c
  - curl -L https://getenvoy.io/cli | bash -s -- -b /usr/local/bin
- - sh
  - -c
  - nohup getenvoy run standard:1.14.1 -- --config-path /etc/envoy/config.yaml > /dev/null
    2>&1 &
`
