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

func BKBinRead(fileName string, readWholeFile bool) (BKBin, error) {
	var result BKBin

	file, err := os.Open(fileName)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hdr := BKBinHeader{}
	err = binary.Read(file, binary.LittleEndian, &hdr)
	if err != nil {
		return result, err
	}
	result.Header = hdr

	toRead := int(hdr.Length)
	if readWholeFile {
		stat, err := file.Stat()
		if err != nil {
			return result, err
		}
		toRead = int(stat.Size()) - 4
	}
	if toRead <= 0 {
		return result, errors.New("Detected wrong length value, may be it is not a BIN file")
	}

	data := make([]byte, toRead)
	if _, err = io.ReadFull(file, data); err != nil {
		return result, err
	}
	result.Data = data

	return result, nil
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
