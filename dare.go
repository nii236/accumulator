package accumulator

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/minio/sio"
	"golang.org/x/crypto/hkdf"
)

type Darer struct {
	MasterKey []byte
}

func NewDarer(masterKeyHex string) (*Darer, error) {
	masterKey, err := hex.DecodeString(masterKeyHex)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode hex key: %v", err)
	}
	return &Darer{masterKey}, nil
}

func (d *Darer) encrypt(inputB []byte) ([]byte, []byte, error) {
	var nonce [32]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, []byte{}, fmt.Errorf("Failed to read random data: %w", err)
	}

	var key [32]byte
	kdf := hkdf.New(sha256.New, d.MasterKey, nonce[:], nil)
	if _, err = io.ReadFull(kdf, key[:]); err != nil {
		return nil, []byte{}, fmt.Errorf("Failed to derive encryption key: %w", err)
	}
	input := bytes.NewReader(inputB)
	output := &bytes.Buffer{}

	if _, err = sio.Encrypt(output, input, sio.Config{Key: key[:]}); err != nil {
		return nil, []byte{}, fmt.Errorf("Failed to encrypt data: %w", err)
	}
	return output.Bytes(), nonce[:], nil
}
func (d *Darer) decrypt(inputB []byte, nonce []byte) ([]byte, error) {
	var key [32]byte
	kdf := hkdf.New(sha256.New, d.MasterKey, nonce, nil)
	_, err := io.ReadFull(kdf, key[:])
	if err != nil {
		return nil, fmt.Errorf("Failed to derive encryption key: %v", err)
	}

	input := bytes.NewReader(inputB)
	output := &bytes.Buffer{}

	_, err = sio.Decrypt(output, input, sio.Config{Key: key[:]})
	if err != nil {
		if _, ok := err.(sio.Error); ok {
			return nil, fmt.Errorf("Malformed encrypted data: %v", err)
		}
		return nil, fmt.Errorf("Failed to decrypt data: %v", err)
	}
	return output.Bytes(), nil
}
