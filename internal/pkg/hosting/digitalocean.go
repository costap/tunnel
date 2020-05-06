package hosting

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"time"
)

type DigitalOcean struct {
	token  string
	region string
	name   string
	size   string
	image  string
	client *godo.Client
}

//NewDOProvider
func NewDOProvider() *DigitalOcean {
	t := "7094b1d59d6354ff71207cdb53de33df85e86b9c1e2a9781284ecdbe9e03cc14"
	return &DigitalOcean{
		token:  t,
		region: "lon1",
		name:   "provision-droplet",
		size:   "s-1vcpu-1gb",
		image:  "ubuntu-18-04-x64",
		client: godo.NewFromToken(t),
	}
}

func (p DigitalOcean) CreateInstance() (string, error) {

	createRequest := &godo.DropletCreateRequest{
		Name:   p.name,
		Region: p.region,
		Size:   p.size,
		Image: godo.DropletCreateImage{
			Slug: p.image,
		},
	}

	ctx := context.TODO()

	newDroplet, _, err := p.client.Droplets.Create(ctx, createRequest)

	if err != nil {
		return "", fmt.Errorf("error creating droplet %w", err)
	}

	id := newDroplet.ID

	fmt.Printf("Created new droplet with id: %v\n", id)

	ip := ""

	for ip == "" {
		time.Sleep(1 * time.Second)
		status, _ip := p.checkStatusIp(id)
		fmt.Printf("droplet %v status: %v\n", id, status)
		ip = _ip
	}

	return ip, nil
}

func (p DigitalOcean) checkStatusIp(id int) (string, string) {

	d, _, err := p.client.Droplets.Get(context.Background(), id)
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
