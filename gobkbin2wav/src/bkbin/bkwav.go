package bkbin

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math"
	"os"
)

var SND_PARTS = [][]byte{
	[]byte{0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40},
	[]byte{0x80, 0xa0, 0xb7, 0xc0, 0xb7, 0xa0, 0x80, 0x5f, 0x48, 0x3f, 0x48, 0x5f, 0x80, 0xb7, 0xb7, 0x80, 0x48, 0x48},
	[]byte{0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40,
		0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40,
		0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40},
	[]byte{0x80, 0x90, 0x9d, 0xa4, 0xa6, 0xa9, 0xa9, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0x80, 0x6f, 0x62, 0x5b, 0x57, 0x56,
		0x55, 0x55, 0x55, 0x55, 0x55, 0x6e, 0x80, 0x9a, 0xa2, 0xa6, 0xa7, 0xa9, 0x80, 0x6f, 0x63, 0x5c, 0x59, 0x59,
		0x80, 0xb7, 0xb7, 0x80, 0x48, 0x48},
	[]byte{0x80, 0x8f, 0xa8, 0xb5, 0xbc, 0xbf, 0xc1, 0xc1, 0xc2, 0xc2, 0xc2, 0xc2, 0xc2, 0xc1, 0xc1, 0xc1, 0xc1, 0xc1,
		0xc0, 0xc0, 0xc0, 0xc0, 0xc0, 0xbf, 0xbf, 0xbe, 0xbe, 0xbe, 0xbd, 0xbc, 0xb2, 0x80, 0x59, 0x48, 0x40, 0x3b,
		0x39, 0x38, 0x37, 0x37, 0x37, 0x37, 0x37, 0x38, 0x38, 0x38, 0x38, 0x39, 0x39, 0x39, 0x39, 0x39, 0x3a, 0x3a,
		0x3a, 0x3a, 0x3a, 0x3b, 0x3b, 0x3b, 0x3b, 0x3c, 0x3e, 0x7a}}

const SIGNAL_RESET = 0
const SIGNAL_SET = 1
const SIGNAL_START_SYNCHRO = 2
const SIGNAL_SYNCHRO = 3
const SIGNAL_END_MARKER = 4

type WavChunkHeader struct {
	ID   [4]uint8
	Size uint32
}

type WavFormat struct {
	ID            [4]uint8
	Size          uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

func writeHeader(writer io.Writer, length uint32) error {
	header := WavChunkHeader{}
	copy(header.ID[:], []uint8("RIFF")[0:4])
	header.Size = length

	if err := binary.Write(writer, binary.LittleEndian, header); err != nil {
		return err
	}

	return binary.Write(writer, binary.LittleEndian, []uint8("WAVE"))
}

func writeObj(writer io.Writer, obj interface{}) error {
	if err := binary.Write(writer, binary.LittleEndian, obj); err != nil {
		return err
	}
	return nil
}

func writeSndByte(target *bytes.Buffer, value uint8) {
	for i := 0; i < 8; i++ {
		writeSndSignal(target, int(value&1), 1)
		value = value >> 1
	}
}

func writeSndShort(target *bytes.Buffer, value uint16) {
	writeSndByte(target, byte(value))
	writeSndByte(target, byte(value>>8))
}

func writeSndArray(target *bytes.Buffer, value []uint8) {
	for _, element := range value {
		writeSndByte(target, element)
	}
}

func writeSndSignal(target *bytes.Buffer, index int, times int) {
	signal := SND_PARTS[index]
	for i := 0; i < times; i++ {
		if _, err := target.Write(signal); err != nil {
			log.Fatal(err)
		}
	}
}

func StringToBK0010(s string) []byte {
	runes := []rune(s)
	result := make([]byte, 0, len(runes))

	for _, r := range runes {
		code := runeToBK0010(r)
		if code != 0 {
			result = append(result, code)
		}
	}

	return result
}

var cyrillicMap = map[rune]uint8{
	'ю': 192, 'а': 193, 'б': 194, 'ц': 195,
	'д': 196, 'е': 197, 'ф': 198, 'г': 199,
	'х': 200, 'и': 201, 'й': 202, 'к': 203,
	'л': 204, 'м': 205, 'н': 206, 'о': 207,
	'п': 208, 'я': 209, 'р': 210, 'с': 211,
	'т': 212, 'у': 213, 'ж': 214, 'в': 215,
	'ь': 216, 'ы': 217, 'з': 218, 'ш': 219,
	'э': 220, 'щ': 221, 'ч': 222, 'ъ': 223,

	'Ю': 224, 'А': 225, 'Б': 226, 'Ц': 227,
	'Д': 228, 'Е': 229, 'Ф': 230, 'Г': 231,
	'Х': 232, 'И': 233, 'Й': 234, 'К': 235,
	'Л': 236, 'М': 237, 'Н': 238, 'О': 239,
	'П': 240, 'Я': 241, 'Р': 242, 'С': 243,
	'Т': 244, 'У': 245, 'Ж': 246, 'В': 247,
	'Ь': 248, 'Ы': 249, 'З': 250, 'Ш': 251,
	'Э': 252, 'Щ': 253, 'Ч': 254, 'Ъ': 255,

	'ё': 197, 'Ё': 229,
}

func runeToBK0010(r rune) byte {
	// ASCII characters (0x00-0x7F)
	if r >= 0x00 && r <= 0x7F {
		return byte(r)
	}

	if code, ok := cyrillicMap[r]; ok {
		return byte(code)
	}

	// Unmapped character - dot
	return 46
}

func writeSndName(target *bytes.Buffer, name string) {
	writtenChars := 0
	nameBk0010 := StringToBK0010(name)
	for i, c := range nameBk0010 {
		if i >= 16 {
			break
		}
		writeSndByte(target, uint8(c))
		writtenChars++
	}
	for writtenChars < 16 {
		writeSndByte(target, uint8(' '))
		writtenChars++
	}
}

func round(f float64) int32 {
	if f < -0.5 {
		return int32(f - 0.5)
	}
	if f > 0.5 {
		return int32(f + 0.5)
	}
	return 0
}

func restrictInByte(a int32) int32 {
	if a < 0 {
		return 0
	}
	if a > 255 {
		return 255
	}
	return a
}

func amplifySnd(data *[]byte) {
	var minDetectedLevel = 256
	var maxDetectedLevel = 0

	for _, e := range *data {
		var ie = int(e)
		if ie < minDetectedLevel {
			minDetectedLevel = ie
		}
		if ie > maxDetectedLevel {
			maxDetectedLevel = ie
		}
	}

	maxDetectedLevel -= 128
	minDetectedLevel -= 128

	var c_max float64 = 128.0 / float64(maxDetectedLevel)
	var c_min float64 = -127.0 / float64(minDetectedLevel)

	var coeff = math.Min(c_max, c_min)

	for i, e := range *data {
		(*data)[i] = byte(restrictInByte(round((float64(e)-128.0)*coeff) + 128))
	}
}

func makeSoundData(bin *BKBin, name string, turbo bool) ([]byte, uint16) {
	buffer := bytes.NewBuffer(make([]byte, 1024*1024))
	buffer.Reset()

	writeSndSignal(buffer, SIGNAL_START_SYNCHRO, 512)
	writeSndSignal(buffer, SIGNAL_SYNCHRO, 1)

	writeSndSignal(buffer, SIGNAL_START_SYNCHRO, 1)
	writeSndSignal(buffer, SIGNAL_SYNCHRO, 1)

	writeSndShort(buffer, bin.Header.Start)
	writeSndShort(buffer, bin.Header.Length)

	writeSndName(buffer, name)

	writeSndSignal(buffer, SIGNAL_START_SYNCHRO, 1)
	writeSndSignal(buffer, SIGNAL_SYNCHRO, 1)

	writeSndArray(buffer, bin.Data)

	var checksum uint16 = CalcChecksum(bin)
	writeSndShort(buffer, checksum)

	writeSndSignal(buffer, SIGNAL_END_MARKER, 1)
	if turbo {
		writeSndSignal(buffer, SIGNAL_START_SYNCHRO, 64)
	} else {
		writeSndSignal(buffer, SIGNAL_START_SYNCHRO, 32)
	}
	writeSndSignal(buffer, SIGNAL_SYNCHRO, 1)

	return buffer.Bytes(), checksum
}

func WriteWavIntoWriter(w io.Writer, name string, turbo bool, amplify bool, bin *BKBin) (uint16, error) {
	var err error = nil
	var checksum uint16 = 0

	sndData, checksum := makeSoundData(bin, name, turbo)
	if amplify {
		amplifySnd(&sndData)
	}

	if err = writeHeader(w, uint32(36+len(sndData))); err != nil {
		return checksum, err
	}

	wavFormat := WavFormat{}
	copy(wavFormat.ID[:], []uint8("fmt ")[0:4])
	wavFormat.Size = 16
	wavFormat.AudioFormat = 1
	wavFormat.NumChannels = 1
	if turbo {
		wavFormat.SampleRate = 22050
		wavFormat.ByteRate = 22050
	} else {
		wavFormat.SampleRate = 11025
		wavFormat.ByteRate = 11025
	}
	wavFormat.BlockAlign = 1
	wavFormat.BitsPerSample = 8

	wavData := WavChunkHeader{}
	copy(wavData.ID[:], []uint8("data")[0:4])
	wavData.Size = uint32(len(sndData))

	if err = writeObj(w, wavFormat); err != nil {
		return checksum, err
	}

	if err = writeObj(w, wavData); err != nil {
		return checksum, err
	}

	if _, err = w.Write(sndData); err != nil {
		return checksum, err
	}

	return checksum, nil
}

func WriteWav(targetFileName string, name string, turbo bool, amplify bool, bin *BKBin) (uint16, error) {
	file, err := os.Create(targetFileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return WriteWavIntoWriter(file, name, turbo, amplify, bin)
}
