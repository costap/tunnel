package do

import "github.com/digitalocean/godo"

type DigitalOcean struct {
	token       string
	client      *godo.Client
}

//NewDOProvider
func NewDOProvider(token string) *DigitalOcean {
	return &DigitalOcean{
		token:  token,
		client: godo.NewFromToken(token),
	}
}

