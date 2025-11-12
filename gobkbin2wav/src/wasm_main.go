//go:build js || wasm

package main

import (
	"bytes"
	"fmt"
	"strconv"
	"syscall/js"

	bkbin "github.com/raydac/bkbin2wav/bkbin"
)

// MakeWavFromBk0010Bin(data []byte, useFileSize bool, amplify bool, turbo bool, tapHeaderName string, addressStart int) map { "data":..., "error":...}
func MakeWavFromBk0010Bin(this js.Value, args []js.Value) interface{} {
	if len(args) < 6 {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "missing arguments",
		})
	}

	src := args[0]

	if src.Type() != js.TypeObject {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "first argument must be Uint8Array",
		})
	}

	consoleOutput := func(s string) {
		js.Global().Get("console").Call("log", s)
	}

	length := src.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, src)

	consoleOutput("Received " + strconv.Itoa(len(data)) + " bytes")

	useFileSize := args[1].Bool()
	amplify := args[2].Bool()
	turbo := args[3].Bool()
	tapeFileName := args[4].String()
	addressStart := args[5].Int()

	bkBin, err := bkbin.BKBinReadFromReader(BytesToReader(data), int64(len(data)), useFileSize)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "error during parse BIN format: " + err.Error(),
		})
	}

	if addressStart < 0 {
		addressStart = int(bkBin.Header.Start)
	}

	if useFileSize {
		consoleOutput(fmt.Sprintf("Detected flag to enforce physical file size (size defined inside of .BIN is %d byte(s), real size is %d byte(s))\n", bkBin.Header.Length, len(bkBin.Data)))
	} else {
		if int(bkBin.Header.Length) != int(length-4) {
			consoleOutput(fmt.Sprintf("Detected different size defined in BIN header (%d != %d)\n", bkBin.Header.Length, length-4))
		}
	}

	if addressStart != int(bkBin.Header.Start) {
		consoleOutput(fmt.Sprintf("The Start address has been changed from %d(&O%o) to %d(&O%o)\n", bkBin.Header.Start, bkBin.Header.Start, addressStart, addressStart))
		bkBin.Header.Start = uint16(addressStart)
	}

	if len(tapeFileName) == 0 || len(tapeFileName) > 16 {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "illegal tape file name: " + tapeFileName,
		})
	}

	var wavBuffer bytes.Buffer
	_, err = bkbin.WriteWavIntoWriter(&wavBuffer, tapeFileName, turbo, amplify, bkBin)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "error during WAV write: " + err.Error(),
		})
	}

	resultWavData := wavBuffer.Bytes()
	consoleOutput("Generated WAV data " + strconv.Itoa(len(resultWavData)) + " bytes")

	uint8Array := js.Global().Get("Uint8Array").New(len(resultWavData))
	js.CopyBytesToJS(uint8Array, resultWavData)

	consoleOutput("Converted into Uint8Array " + strconv.Itoa(uint8Array.Get("length").Int()) + " bytes")

	return js.ValueOf(map[string]interface{}{
		"data":  uint8Array,
		"error": nil,
	})
}

func main() {
	js.Global().Set("MakeWavFromBk0010Bin", js.FuncOf(MakeWavFromBk0010Bin))
	println("WASM Go runtime started")
	select {}
}
