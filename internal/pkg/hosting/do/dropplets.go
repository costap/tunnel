package do

import (
	"context"
	"fmt"
	"github.com/costap/tunnel/internal/pkg/envoy"
	"github.com/digitalocean/godo"
	"log"
	"time"
)

const DefaultSize = "s-1vcpu-1gb"

type Proxy struct {
	SourcePort, TargetPort int
}

type HostConfig struct {
	region      string
	name        string
	size        string
	image       string
	sshKeysPath string
	sshKeyName  string
	proxies     []Proxy
}

func NewHostConfig(region, name, image, sshKeysPath, sshKeyName string) *HostConfig {
	return &HostConfig{
		region: region,
		name: name,
		image: image,
		sshKeysPath: sshKeysPath,
		sshKeyName: sshKeyName,
		size: DefaultSize,
		proxies: []Proxy{},
	}
}

func (c *HostConfig) SetHostName(name string){
	c.name = name
}

func (c *HostConfig) AddProxy(sourcePort, targetPort int) *HostConfig {
	c.proxies = append(c.proxies, Proxy{
		SourcePort: sourcePort,
		TargetPort: targetPort,
	})
	return c
}

func (do DigitalOcean) CreateInstance(config *HostConfig) (string, error) {

	key, err := do.GetKeyOrCreate(config.sshKeysPath, config.sshKeyName)
	if err != nil {
		return "", fmt.Errorf("error getting ssh key: %w", err)
	}
	log.Printf("Using key: %v\n", key)
	c := envoy.NewConfig()
	for i, p := range config.proxies {
		c.AddTCPProxy(envoy.TCPProxy{
			Name:         fmt.Sprintf("proxy%v", i),
			ListenerPort: p.SourcePort,
			ClusterPort:  p.TargetPort,
		})
	}
	createRequest := &godo.DropletCreateRequest{
		Name:   config.name,
		Region: config.region,
		Size:   config.size,
		Image: godo.DropletCreateImage{
			Slug: config.image,
		},
		SSHKeys: []godo.DropletCreateSSHKey{{
			Fingerprint: key.Fingerprint,
		}},
		UserData: c.CloudConfigYaml(),
	}

	ctx := context.TODO()

	newDroplet, _, err := do.client.Droplets.Create(ctx, createRequest)

	if err != nil {
		return "", fmt.Errorf("error creating droplet %w", err)
	}

	id := newDroplet.ID

	fmt.Printf("Created new droplet with id: %v\n", id)

	ip := ""

	for ip == "" {
		time.Sleep(1 * time.Second)
		status, _ip := do.checkStatusIp(id)
		fmt.Printf("waiting for ip for droplet %v status: %v...\n", id, status)
		ip = _ip
	}

	return ip, nil
}

func (do DigitalOcean) checkStatusIp(id int) (string, string) {

	d, _, err := do.client.Droplets.Get(context.Background(), id)
	if err != nil {
		fmt.Printf("error checking status: %v\n", err)
		return "Unknown", ""
	}

	ip, err := d.PublicIPv4()
	if err != nil {
		fmt.Printf("error checking ip: %v\n", err)
		return d.Status, ""
	}
	return d.Status, ip
}
