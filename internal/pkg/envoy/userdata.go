package envoy

import (
	"gopkg.in/yaml.v2"
	"log"
)

type CloudConfig struct {
	WriteFiles []struct{
		Content string `yaml:"content"`
		Path string `yaml:"path"`
	} `yaml:"write_files"`
	RunCMD [][]string `yaml:"runcmd"`
}

func (c *Config) CloudConfigYaml() string {
	t := CloudConfig{
		WriteFiles: []struct{
			Content string `yaml:"content"`
			Path string `yaml:"path"`
		}{
			{Content: c.ToYaml(), Path: "/etc/envoy/config.yaml"},
		},
		RunCMD: [][]string{
			{"sh","-c","curl -L https://getenvoy.io/cli | bash -s -- -b /usr/local/bin"},
			{"sh","-c","getenvoy run standard:1.14.1 -- --config-path /etc/envoy/config.yaml"},
		},
	}
	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return "#cloud-config\n" + string(d)
}


