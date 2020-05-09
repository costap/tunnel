package do

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/costap/tunnel/internal/pkg/keys"
	"github.com/digitalocean/godo"
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
	return do.checkGetKeyOrCreatePerPage(localKeysPath, localKeyName, 1, 100)
}

func (do *DigitalOcean) checkGetKeyOrCreatePerPage(localKeysPath, localKeyName string, page, pageSize int) (*godo.Key, error) {

	lk, err := do.ListKeys(page, pageSize)
	if err != nil {
		log.Printf("error listing keys from DO, will try create: %v", err)
		return do.CreateKey(localKeysPath, localKeyName, "tunnel")
	}

	if lk == nil || len(lk) == 0 {
		log.Println("no keys found in DO, will create")
		return do.CreateKey(localKeysPath, localKeyName, "tunnel")
	}

	if _, err := os.Stat(filepath.Join(localKeysPath, localKeyName+".pub")); os.IsNotExist(err) {
		log.Print("local key not found, will create")
		return do.CreateKey(localKeysPath, localKeyName, "tunnel")
	}

	for _, k := range lk {
		if lpk, _ := keys.ReadPublicKey(localKeysPath, localKeyName); k.PublicKey == lpk {
			log.Print("found matching key in DO, will use")
			return &k, nil
		}
	}

	if len(lk) == pageSize {
		log.Print("no matching keys found in page, will check next page")
		return do.checkGetKeyOrCreatePerPage(localKeysPath, localKeyName, page+1, pageSize)
	}

	log.Print("no key in DO matching local key, will create")
	return do.CreateKey(localKeysPath, localKeyName, "tunnel")
}

func (do *DigitalOcean) CreateKey(localKeysPath, localKeyName, name string) (*godo.Key, error) {

	if _, err := os.Stat(filepath.Join(localKeysPath, localKeyName+".pub")); os.IsNotExist(err) {
		keys.GenerateKeyPair(localKeysPath, localKeyName)
	}

	pkBytes, err := ioutil.ReadFile(filepath.Join(localKeysPath, localKeyName+".pub"))
	if err != nil {
		return nil, fmt.Errorf("error reading public key: %w", err)
	}

	key, _, err := do.client.Keys.Create(context.Background(), &godo.KeyCreateRequest{
		Name:      name + "_" + localKeyName,
		PublicKey: string(pkBytes),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating DO key: %w", err)
	}

	return key, nil
}
