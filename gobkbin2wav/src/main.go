package main

import (
	"bytes"
	"io"
)

func BytesToReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}
