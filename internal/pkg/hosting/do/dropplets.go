package do

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"time"
)

func (do DigitalOcean) CreateInstance() (string, error) {

	key, err := do.GetKeyOrCreate(do.sshKeysPath, do.sshKeyName)
	if err != nil {
		return "", fmt.Errorf("error getting ssh key: %w", err)
	}
	createRequest := &godo.DropletCreateRequest{
		Name:   do.name,
		Region: do.region,
		Size:   do.size,
		Image: godo.DropletCreateImage{
			Slug: do.image,
		},
		SSHKeys: []godo.DropletCreateSSHKey{{
			Fingerprint: key.Fingerprint,
		}},
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
		fmt.Printf("droplet %v status: %v\n", id, status)
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
