package main

import (
	"github.com/3dev/goks"
	"path/filepath"
)

func openKeyStore(filename string, passkey string) (*goks.KeyStore, error) {

	ext := filepath.Ext(filename)

	// If there's no extension, add ".goks"
	if ext == "" {
		filename += ".goks"
	}

	ks, err := goks.Open(filename, passkey)
	if err != nil {
		return nil, err
	}

	return ks, nil
}
