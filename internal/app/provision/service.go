package provision

import (
	"fmt"
	"github.com/costap/tunnel/internal/pkg/hosting"
)

type ProvisionService struct {
}

func NewProvisionService() *ProvisionService {
	return &ProvisionService{}
}

func (*ProvisionService) CreateHosting() {
	p := hosting.NewDOProvider()

	ip, err := p.CreateInstance()

	if err != nil {
		fmt.Printf("error creating host %v\n", err)
		return
	}

	fmt.Printf("host created with ip %v\n", ip)
}
