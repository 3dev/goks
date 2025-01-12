package goKeyStore

import (
	"bufio"
	"bytes"
	"errors"
	"goKeyStore/crypt"
	"goKeyStore/file"
	"os"
)

type (
	KeyStore struct {
		fileReader *bufio.Reader
		fileWriter *bufio.Writer
		passkey    string
		fh         *file.FileHeader
	}
)

var (
	ErrUnableToValidateCheck = errors.New("unable to validate check digits")
	ErrCheckDigitFailed      = errors.New("check digits failed")
	ErrFormattingKeyStore    = errors.New("formatting keystore failed")
)

func New(filename string, passkey string) (*KeyStore, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	ks := &KeyStore{bufio.NewReader(f), bufio.NewWriter(f), passkey, &file.FileHeader{}}
	err = ks.formatKeyStore(passkey)
	if err != nil {
		return nil, err
	}

	return ks, nil
}

func Open(filename string, passkey string) (*KeyStore, error) {
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	kStore := &KeyStore{bufio.NewReader(f), bufio.NewWriter(f), passkey, nil}
	if err = kStore.validatePasskey(passkey); err != nil {
		return nil, err
	}

	return kStore, nil
}

func (ks *KeyStore) Close() error {
	return ks.fileWriter.Flush()
}

func (ks *KeyStore) validatePasskey(passkey string) error {

	fileCheckDigit := make([]byte, 4)
	_, err := ks.fileReader.Read(fileCheckDigit)
	if err != nil {
		return err
	}

	key, _ := crypt.PadAESKey([]byte(passkey))
	checkDigit, err := crypt.EncryptAESECB(key, bytes.Repeat([]byte{0}, 16))
	if err != nil {
		return ErrUnableToValidateCheck
	}

	if !bytes.Equal(fileCheckDigit[:4], checkDigit[:4]) {
		return ErrCheckDigitFailed
	}

	headerBytes := make([]byte, file.HeaderSize-4)
	_, err = ks.fileReader.Read(headerBytes)
	if err != nil {
		return err
	}

	err = ks.fh.Decode(append(checkDigit, headerBytes...))
	if err != nil {
		return err
	}

	return nil
}

func (ks *KeyStore) formatKeyStore(passkey string) error {

	key, _ := crypt.PadAESKey([]byte(passkey))
	checkDigit, err := crypt.EncryptAESECB(key, bytes.Repeat([]byte{0}, 16))
	if err != nil {
		return ErrFormattingKeyStore
	}

	fh := file.FileHeader{}
	copy(fh.CheckDigit[:], checkDigit)
	_, err = ks.fileWriter.Write(fh.Bytes())
	if err != nil {
		return err
	}

	return ks.fileWriter.Flush()
}

func (ks *KeyStore) readHeader() error {

	headerBytes := make([]byte, file.HeaderSize)
	_, err := ks.fileReader.Read(headerBytes)
	if err != nil {
		return err
	}
	return ks.fh.Decode(headerBytes)
}
