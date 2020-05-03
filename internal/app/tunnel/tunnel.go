package tunnel

import (
	"fmt"
	"github.com/costap/tunnel/internal/pkg/hosting"
)

type Tunnel struct {

}

func NewTunnel() *Tunnel {
	return &Tunnel{}
}

func (*Tunnel) CreateHosting(){
	p := hosting.NewDOProvider()

	ip, err := p.CreateInstance()

	if err != nil {
		fmt.Printf("error creating host %v\n", err)
		return
	}

	fmt.Printf("host created with ip %v\n", ip)
}