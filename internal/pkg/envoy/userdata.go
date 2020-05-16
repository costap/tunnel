package envoy

import (
	"log"

	"gopkg.in/yaml.v2"
)

type CloudConfig struct {
	WriteFiles []struct {
		Content string `yaml:"content"`
		Path    string `yaml:"path"`
	} `yaml:"write_files"`
	RunCMD [][]string `yaml:"runcmd"`
}

func (c *Config) CloudConfigYaml() string {
	t := CloudConfig{
		WriteFiles: []struct {
			Content string `yaml:"content"`
			Path    string `yaml:"path"`
		}{
			{Content: openfiles, Path: "/etc/security/limits.conf"},
			{Content: sshdConfig, Path: "/etc/ssh/sshd_config"},
			{Content: c.ToYaml(), Path: "/etc/envoy/config.yaml"},
		},
		RunCMD: [][]string{
			{"sh", "-c", "curl -L https://getenvoy.io/cli | bash -s -- -b /usr/local/bin"},
			// TODO: Set envoy with systemctl
			{"sh", "-c", "nohup getenvoy run standard:1.14.1 -- --config-path /etc/envoy/config.yaml > /dev/null 2>&1 &"},
		},
	}
	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return "#cloud-config\n" + string(d)
}

const openfiles = `*         hard    nofile      500000
*         soft    nofile      500000
root      hard    nofile      500000
root      soft    nofile      500000
`

const sshdConfig = "PermitRootLogin yes\n" +
					"PasswordAuthentication no\n" +
					"ChallengeResponseAuthentication no\n" +
					"UsePAM yes\n" +
					"X11Forwarding yes\n" +
					"PrintMotd no\n" +
					"ClientAliveInterval 120\n" +
					"ClientAliveCountMax 720\n" +
					"AcceptEnv LANG LC_*\n" +
					"Subsystem	sftp	/usr/lib/openssh/sftp-server"
