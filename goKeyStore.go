package goKeyStore

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/3dev/goKeyStore/crypt"
	"github.com/3dev/goKeyStore/file"
	"io"
	"os"
	"path/filepath"
	"slices"
)

type (
	KeyStore struct {
		keyStoreFile *os.File
		passkey      string
		fileHeader   *file.FileHeader
		itemCount    int
	}
)

const (
	TblContentSize     = 1024
	FirstContentPos    = file.HeaderSize
	TblContentStartPos = 4
	TblContentItemSize = 45
)

var (
	ErrUnableToValidateCheck = errors.New("unable to validate check digits")
	ErrCheckDigitFailed      = errors.New("check digits failed")
	ErrFormattingKeyStore    = errors.New("formatting keystore failed")
	ErrKeyStoreFull          = errors.New("keystore full")
	ErrKeyTooLarge           = errors.New("key is too large")
	ErrNotFound              = errors.New("key not found")
	ErrDuplicateKey          = errors.New("duplicate key")
)

func New(filename string, passkey string) (*KeyStore, error) {

	ext := filepath.Ext(filename)

	// If there's no extension, add ".gks"
	if ext == "" {
		filename += ".goks"
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	ks := &KeyStore{
		keyStoreFile: f,
		passkey:      passkey,
		fileHeader:   nil,
		itemCount:    0,
	}

	fh, err := ks.formatKeyStore(passkey)
	ks.fileHeader = fh
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

	kStore := &KeyStore{
		keyStoreFile: f,
		passkey:      passkey,
		fileHeader:   &file.FileHeader{},
		itemCount:    0,
	}

	if err = kStore.validatePasskey(passkey); err != nil {
		return nil, err
	}

	return kStore, nil
}

func (ks *KeyStore) Close() error {
	return ks.keyStoreFile.Close()
}

func (ks *KeyStore) validatePasskey(passkey string) error {

	fileCheckDigit := make([]byte, 4)
	_, err := ks.keyStoreFile.Read(fileCheckDigit)
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
	_, err = ks.keyStoreFile.Read(headerBytes)
	if err != nil {
		return err
	}

	count, err := ks.fileHeader.Decode(append(fileCheckDigit, headerBytes...))
	if err != nil {
		return err
	}

	ks.itemCount = count

	return nil
}

func (ks *KeyStore) formatKeyStore(passkey string) (*file.FileHeader, error) {

	key, _ := crypt.PadAESKey([]byte(passkey))
	checkDigit, err := crypt.EncryptAESECB(key, bytes.Repeat([]byte{0}, 16))
	if err != nil {
		return nil, ErrFormattingKeyStore
	}

	fh := file.FileHeader{}
	copy(fh.CheckDigit[:], checkDigit)
	_, err = ks.keyStoreFile.Write(fh.Bytes())
	if err != nil {
		return nil, err
	}

	err = ks.keyStoreFile.Sync()
	if err != nil {
		return nil, err
	}

	return &fh, nil
}

func (ks *KeyStore) Count() int {
	return ks.itemCount
}

func (ks *KeyStore) Keys() []string {

	keys := make([]string, ks.itemCount)
	iKeys := 0
	for i := 0; i < TblContentSize; i++ {
		if ks.fileHeader.Index[i].Available > 0 {
			keys[iKeys] = string(bytes.TrimRight(ks.fileHeader.Index[i].Key[:], string([]byte{0})))
			iKeys++
		}
	}
	return keys
}

func bytesPad(s string, maxBytes int) ([]byte, error) {

	if len(s) > maxBytes {
		return nil, errors.New("string longer than max bytes")
	}

	return append([]byte(s), bytes.Repeat([]byte{0}, maxBytes-len(s))...), nil
}

func (ks *KeyStore) Put(key string, data []byte) error {

	if ks.itemCount >= TblContentSize {
		return ErrKeyStoreFull
	}

	if slices.Contains(ks.Keys(), key) {
		return ErrDuplicateKey
	}

	//find a free index
	filePos := uint32(FirstContentPos)
	for i := 0; i < TblContentSize; i++ {
		l := binary.BigEndian.Uint32(ks.fileHeader.Index[i].AllocatedLength[:])
		if ks.fileHeader.Index[i].Available < 1 {

			if l < 1 { //it is unused and the next available spot
				k, err := bytesPad(key, 32)
				if err != nil {
					return ErrKeyTooLarge
				}
				copy(ks.fileHeader.Index[i].Key[:], k)
				ks.fileHeader.Index[i].Available = 1
				binary.BigEndian.PutUint32(ks.fileHeader.Index[i].DataLength[:], uint32(len(data)))
				binary.BigEndian.PutUint32(ks.fileHeader.Index[i].AllocatedLength[:], uint32(len(data)))
				binary.BigEndian.PutUint32(ks.fileHeader.Index[i].Location[:], filePos)
				_, err = ks.keyStoreFile.Seek(int64(filePos), io.SeekStart)
				if err != nil {
					return err
				}
				_, err = ks.keyStoreFile.Write(data)
				if err != nil {
					return err
				}

				//write the index
				_, err = ks.keyStoreFile.Seek(int64(TblContentStartPos+(i*TblContentItemSize)), io.SeekStart)
				if err != nil {
					return err
				}
				_, err = ks.keyStoreFile.Write(ks.fileHeader.Index[i].Bytes())
				if err != nil {
					return err
				}
				err = ks.keyStoreFile.Sync()
				if err != nil {
					return err
				}

				ks.itemCount++
				break
			}

			//it was previously used but the item was marked deleted
			if uint32(len(data)) < l {
				//we can use it
				k, err := bytesPad(key, 32)
				if err != nil {
					return ErrKeyTooLarge
				}
				copy(ks.fileHeader.Index[i].Key[:], k)
				ks.fileHeader.Index[i].Available = 1
				binary.BigEndian.PutUint32(ks.fileHeader.Index[i].DataLength[:], uint32(len(data)))
				//binary.BigEndian.PutUint32(ks.fileHeader.Index[i].Location[:], filePos)
				_, err = ks.keyStoreFile.Seek(int64(binary.BigEndian.Uint32(ks.fileHeader.Index[i].Location[:])), io.SeekStart)
				if err != nil {
					return err
				}
				_, err = ks.keyStoreFile.Write(data)
				if err != nil {
					return err
				}

				//write the index
				_, err = ks.keyStoreFile.Seek(int64(TblContentStartPos+(i*TblContentItemSize)), io.SeekStart)
				if err != nil {
					return err
				}
				_, err = ks.keyStoreFile.Write(ks.fileHeader.Index[i].Bytes())
				if err != nil {
					return err
				}
				err = ks.keyStoreFile.Sync()
				if err != nil {
					return err
				}

				break
			}
		}

		filePos += l
	}

	return nil
}

func (ks *KeyStore) Delete(key string) error {

	for i := 0; i < TblContentSize; i++ {
		if ks.fileHeader.Index[i].Available > 0 {
			ksKey := string(bytes.TrimRight(ks.fileHeader.Index[i].Key[:], string([]byte{0})))
			if ksKey == key {
				ks.fileHeader.Index[i].Available = 0
				_, err := ks.keyStoreFile.Seek(int64(TblContentStartPos+(i*TblContentItemSize)), io.SeekStart)
				if err != nil {
					return err
				}
				_, err = ks.keyStoreFile.Write(ks.fileHeader.Index[i].Bytes())
				if err != nil {
					return err
				}
				err = ks.keyStoreFile.Sync()
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return ErrNotFound
}

func (ks *KeyStore) Get(key string) ([]byte, error) {

	for i := 0; i < TblContentSize; i++ {
		if ks.fileHeader.Index[i].Available > 0 {
			ksKey := string(bytes.TrimRight(ks.fileHeader.Index[i].Key[:], string([]byte{0})))
			if ksKey == key {
				dataLen := binary.BigEndian.Uint32(ks.fileHeader.Index[i].DataLength[:])
				data := make([]byte, dataLen)
				_, err := ks.keyStoreFile.Seek(int64(binary.BigEndian.Uint32(ks.fileHeader.Index[i].Location[:])), io.SeekStart)
				if err != nil {
					return nil, err
				}
				_, err = ks.keyStoreFile.Read(data)
				if err != nil {
					return nil, err
				}

				return data, nil
			}
		}
	}

	return nil, ErrNotFound
}

func (ks *KeyStore) KeyInfo(key string) (file.TableOfContent, error) {

	for i := 0; i < TblContentSize; i++ {
		if ks.fileHeader.Index[i].Available > 0 {
			ksKey := string(bytes.TrimRight(ks.fileHeader.Index[i].Key[:], string([]byte{0})))
			if ksKey == key {
				return ks.fileHeader.Index[i], nil
			}
		}
	}

	return file.TableOfContent{}, ErrNotFound
}

func (ks *KeyStore) Compact() error {

	return nil
}
