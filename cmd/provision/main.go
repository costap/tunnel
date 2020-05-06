package main

import (
	"github.com/costap/tunnel/internal/app/provision"
)

func main() {
	t := provision.NewProvisionService()
	t.CreateHosting()
}
