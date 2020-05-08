package do

import "github.com/digitalocean/godo"

type DigitalOcean struct {
	token       string
	region      string
	name        string
	size        string
	image       string
	sshKeysPath string
	sshKeyName  string
	client      *godo.Client
}

//NewDOProvider
func NewDOProvider(token, region, image, sshKeysPath, sshKeyName string) *DigitalOcean {
	// t := "7094b1d59d6354ff71207cdb53de33df85e86b9c1e2a9781284ecdbe9e03cc14"
	return &DigitalOcean{
		token:  token,
		region: region,
		name:   "tunnelctl-droplet",
		size:   "s-1vcpu-1gb",
		image:  image,
		sshKeysPath: sshKeysPath,
		sshKeyName: sshKeyName,
		client: godo.NewFromToken(token),
	}
}
