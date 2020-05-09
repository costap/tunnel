package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GenerateKeyPair(path, name string) error {
	reader := rand.Reader
	bitSize := 4096

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return fmt.Errorf("error generating rsa key %w", err)
	}

	return SaveKeyPair(key, path, name)
}

func SaveKeyPair(key *rsa.PrivateKey, path, name string) error {

	if err := savePEMKey(filepath.Join(path, name), key); err != nil {
		return fmt.Errorf("error saving pem private key: %w", err)
	}

	pubBytes, err := generatePublicKey(&key.PublicKey)
	if err != nil {
		return fmt.Errorf("error generatin public key bytes: %w", err)
	}

	if err := writeKeyToFile(pubBytes, filepath.Join(path, name + ".pub")); err != nil {
		return fmt.Errorf("error saving pem public key: %w", err)
	}

	return nil
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}

func savePEMKey(fileName string, key *rsa.PrivateKey) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file for pem key: %w", err)
	}
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		return fmt.Errorf("error encoding private key to file: %w", err)
	}

	return nil
}

func ReadPublicKey(path, name string) (string, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(path, name + ".pub"))
	if err != nil {
		return "", fmt.Errorf("error reading public key %w", err)
	}
	return string(bytes), nil
}
