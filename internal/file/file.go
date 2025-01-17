package file

import (
	"bytes"
)

const (
	HeaderSize     = (45 * 1024) + 4
	TblContentSize = 1024
)

type (
	TableOfContent struct {
		Available       byte
		Key             [32]byte
		DataLength      [4]byte
		AllocatedLength [4]byte
		Location        [4]byte
	}

	FileHeader struct {
		CheckDigit [4]byte
		Index      [1024]TableOfContent
	}
)

func (fIdx *TableOfContent) Bytes() []byte {

	buff := bytes.Buffer{}
	buff.WriteByte(fIdx.Available)
	buff.Write(fIdx.Key[:])
	buff.Write(fIdx.DataLength[:])
	buff.Write(fIdx.AllocatedLength[:])
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

func (fIdx *TableOfContent) Decode(rd *bytes.Reader) error {

	var err error

	fIdx.Available, err = rd.ReadByte()
	if err != nil {
		return err
	}
	_, err = rd.Read(fIdx.Key[:])
	if err != nil {
		return err
	}
	_, err = rd.Read(fIdx.DataLength[:])
	if err != nil {
		return err
	}
	_, err = rd.Read(fIdx.AllocatedLength[:])
	if err != nil {
		return err
	}
	_, err = rd.Read(fIdx.Location[:])
	if err != nil {
		return err
	}

	return nil
}

func (fHdr *FileHeader) Decode(data []byte) (int, error) {

	var err error

	buff := bytes.NewReader(data)
	_, err = buff.Read(fHdr.CheckDigit[:])
	if err != nil {
		return 0, err
	}

	c := 0
	for i := 0; i < TblContentSize; i++ {
		err = fHdr.Index[i].Decode(buff)
		if err != nil {
			return 0, err
		}
		if fHdr.Index[i].Available > 0 {
			c++
		}
	}

	return c, nil
}
