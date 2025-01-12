package file

import (
	"bytes"
)

const (
	HeaderSize = (41 * 1024) + 4
)

type (
	FileIndex struct {
		Available byte
		Key       [32]byte
		Length    [4]byte
		Location  [4]byte
	}

	FileHeader struct {
		CheckDigit [4]byte
		Index      [1024]FileIndex
	}
)

func (fIdx *FileIndex) Bytes() []byte {

	buff := bytes.Buffer{}
	buff.WriteByte(fIdx.Available)
	buff.Write(fIdx.Key[:])
	buff.Write(fIdx.Length[:])
	buff.Write(fIdx.Location[:])

	return buff.Bytes()
}

func (fHdr *FileHeader) Bytes() []byte {

	buff := bytes.Buffer{}
	buff.Write(fHdr.CheckDigit[:])
	for i := 0; i < 1024; i++ {
		buff.Write(fHdr.Index[i].Bytes())
	}

	return buff.Bytes()
}

func (fIdx *FileIndex) Decode(data []byte) error {

	var err error

	buff := bytes.NewReader(data)
	fIdx.Available, err = buff.ReadByte()
	if err != nil {
		return err
	}
	_, err = buff.Read(fIdx.Key[:])
	if err != nil {
		return err
	}
	_, err = buff.Read(fIdx.Length[:])
	if err != nil {
		return err
	}
	_, err = buff.Read(fIdx.Location[:])
	if err != nil {
		return err
	}

	return nil
}

func (fHdr *FileHeader) Decode(data []byte) error {

	return nil
}
