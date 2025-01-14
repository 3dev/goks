package main

import (
	"github.com/3dev/goKeyStore"
	"path/filepath"
)

func openKeyStore(filename string, passkey string) (*goKeyStore.KeyStore, error) {

	ext := filepath.Ext(filename)

	// If there's no extension, add ".goks"
	if ext == "" {
		filename += ".goks"
	}

	ks, err := goKeyStore.Open(filename, passkey)
	if err != nil {
		return nil, err
	}

	return ks, nil
}
