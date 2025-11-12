package bkbin

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

type BKBinHeader struct {
	Start  uint16
	Length uint16
}

type BKBin struct {
	Header BKBinHeader
	Data   []uint8
}

func BKBinReadFromReader(r io.Reader, length int64, readWholeFile bool) (*BKBin, error) {
	hdr := BKBinHeader{}
	err := binary.Read(r, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, err
	}

	var result BKBin
	result.Header = hdr

	toRead := int(hdr.Length)
	if readWholeFile {
		toRead = int(length) - 4
	}

	if toRead <= 0 {
		return nil, errors.New("Calculated unexpected data length, may be not a BIN format")
	}

	data := make([]byte, toRead)
	if _, err = io.ReadFull(r, data); err != nil {
		return nil, err
	}
	result.Data = data

	return &result, nil
}

func BKBinRead(fileName string, readWholeFile bool) (*BKBin, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := stat.Size()
	return BKBinReadFromReader(file, size, readWholeFile)
}

func CalcChecksum(bkbin *BKBin) uint16 {
	var sum uint32 = 0

	for _, v := range (*bkbin).Data {
		sum = sum + uint32(v)
		if sum > 0xFFFF {
			sum = (sum & 0xFFFF) + 1
		}
	}
	return uint16(sum)
}
