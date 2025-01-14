package crypt

import (
	"bytes"
	"crypto/aes"
	"errors"
	"fmt"
)

func PadAESKey(key []byte) ([]byte, error) {
	sz := len(key)
	if 16 == sz {
		return key, nil
	} else if sz > 16 {
		return key[:16], nil
	}

	// sz is < 16
	v := key[0] - key[sz-1]
	return append(key, bytes.Repeat([]byte{v}, 16-sz)...), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Remove PKCS#7 padding from plaintext
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}

func EncryptAESECB(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	plaintext = pkcs7Pad(plaintext, block.BlockSize())

	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += block.BlockSize() {
		block.Encrypt(ciphertext[i:i+block.BlockSize()], plaintext[i:i+block.BlockSize()])
	}

	return ciphertext, nil
}

// Decrypt data using AES-ECB
func DecryptAESECB(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, errors.New("file content may be corrupted")
	}

	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += block.BlockSize() {
		block.Decrypt(plaintext[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
	}

	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return nil, errors.New("padding error, file content may be corrupted")
	}

	return plaintext, nil
}
