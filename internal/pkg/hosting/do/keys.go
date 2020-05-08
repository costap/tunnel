package do

import (
	"context"
	"fmt"
	"github.com/costap/tunnel/internal/pkg/keys"
	"github.com/digitalocean/godo"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func (do *DigitalOcean) ListKeys(page, perPage int) ([]godo.Key, error) {
	keys, _, err := do.client.Keys.List(context.Background(), &godo.ListOptions{
		Page:    page,
		PerPage: perPage,
	})

	if err != nil {
		return nil, fmt.Errorf("error listing keys from digitalocean %w", err)
	}
	return keys, nil
}

// GetKeyOrCreate will get a key pair from DO or create one if one doesn't exist
func (do *DigitalOcean) GetKeyOrCreate(localKeysPath, localKeyName string) (*godo.Key, error) {

	keys, err := do.ListKeys(1, 1)
	if err != nil {
		log.Printf("error listing keys from DO, will try create: %v", err)
		return do.CreateKey(localKeysPath, localKeyName,"tunnel")
	}

	if keys == nil || len(keys) == 0 {
		return do.CreateKey(localKeysPath, localKeyName,"tunnel")
	}

	return &keys[0], nil
}

func (do *DigitalOcean) CreateKey(localKeysPath, localKeyName, name string) (*godo.Key, error) {

	if _, err := os.Stat(filepath.Join(localKeysPath, localKeyName + ".pub")); os.IsNotExist(err) {
		keys.GenerateKeyPair(localKeysPath, localKeyName)
	}

	pkBytes, err := ioutil.ReadFile(filepath.Join(localKeysPath, localKeyName + ".pub"))
	if err != nil {
		return nil, fmt.Errorf("error reading public key: %w", err)
	}

	key, _, err := do.client.Keys.Create(context.Background(), &godo.KeyCreateRequest{
		Name:      name,
		PublicKey: string(pkBytes),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating DO key: %w", err)
	}

	return key, nil
}
