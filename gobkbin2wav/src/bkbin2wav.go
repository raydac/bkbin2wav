package main

import (
	"flag"
	"fmt"
	bkbin "github.com/raydac/bkbin2wav/bkbin"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const __AUTHOR__ = "Igor Maznitsa (https://www.igormaznitsa.com)"
const __VERSION__ = "1.0.5"
const __PROJECTURI__ = "https://github.com/raydac/bkbin2wav"

var flagUseFileSize bool
var flagAmplify bool
var flagTurboMode bool
var fileInName string
var fileOutName string
var tapHeaderName string
var addressStart int

func header() {
	fmt.Printf(`
  BKBIN2WAV allows the conversion of .BIN snapshots (for BK-0010(01) emulators) into WAV format.
  Project page : %s
        Author : %s
       Version : %s
  It is a converter that transforms .BIN files (a snapshot format for BK-0010(01) emulators) into WAV sound files compatible with the real BK-0010 TAP reader.
`, __PROJECTURI__, __AUTHOR__, __VERSION__)
}

func init() {
	flag.BoolVar(&flagUseFileSize, "f", false, "use physical file size instead of BIN header value")
	flag.BoolVar(&flagAmplify, "a", false, "amplify audio signal")
	flag.BoolVar(&flagTurboMode, "t", false, "turn on the \"turbo\" mode")
	flag.StringVar(&fileInName, "i", "", "source BIN file")
	flag.StringVar(&fileOutName, "o", "", "target WAV file")
	flag.StringVar(&tapHeaderName, "n", "", "name to be for TAP header (must be less or equals 16 chars)")
	flag.IntVar(&addressStart, "s", -1, "start address")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func assertParameters() os.FileInfo {
	if len(fileInName) == 0 {
		fmt.Fprintf(os.Stdout, `

Examples:

                conversion into WAV : %[1]s -i someBkFile.BIN
  conversion into WAV with new name : %[1]s -i someBkFile.BIN -o wavFile.wav
              conversion into TURBO : %[1]s -t -i someBkFile.BIN
conversion into TURBO with new name : %[1]s -t -i someBkFile.BIN -o wavFile.wav

`, os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(tapHeaderName) > 16 {
		fmt.Fprintln(os.Stderr, "Too long name for TAP file header, must be less or equals 16 chars")
		os.Exit(1)
	}

	if addressStart > 0xFFFF {
		fmt.Fprintln(os.Stderr, "Wrong start address")
		os.Exit(1)
	}

	var inFileInfo, err = os.Stat(fileInName)
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "File doesn't exist : ", fileInName)
		os.Exit(1)
	}
	if inFileInfo.IsDir() {
		fmt.Fprintln(os.Stderr, "Input file must not be folder : ", fileInName)
		os.Exit(1)
	}

	if inFileInfo.Size() <= 4 {
		fmt.Fprintln(os.Stderr, "Input file is too small to be BIN file : ", fileInName)
		os.Exit(1)
	}

	return inFileInfo
}

func asOnOff(flag bool) string {
	if flag {
		return "ON"
	} else {
		return "OFF"
	}
}

func extractName(filename string) string {
	basename := path.Base(filename)
	return strings.ToUpper(strings.TrimSuffix(basename, filepath.Ext(basename)))
}

func main() {
	header()

	flag.Parse()

	var srcFileInfo os.FileInfo = assertParameters()

	binFile, err := bkbin.BKBinRead(fileInName, flagUseFileSize)
	if err != nil {
		log.Fatal(err)
	}

	if len(fileOutName) == 0 {
		fileOutName = filepath.Dir(fileInName) + string(os.PathSeparator) + extractName(fileInName) + ".wav"
	}

	if addressStart < 0 {
		addressStart = int(binFile.Header.Start)
	}

	if flagUseFileSize {
		fmt.Printf("Detected flag to enforce physical file size (size defined inside of .BIN is %d byte(s), real size is %d byte(s))\n", binFile.Header.Length, len(binFile.Data))
	} else {
		if int(binFile.Header.Length) != int(srcFileInfo.Size()-4) {
			fmt.Printf("Warning! Detected different size defined in BIN header, use -f to use file size instead of header size (%d != %d)\n", binFile.Header.Length, srcFileInfo.Size()-4)
		}
	}

	if addressStart != int(binFile.Header.Start) {
		fmt.Printf("Warning! The Start address has been changed from %d(&O%o) to %d(&O%o)\n", binFile.Header.Start, binFile.Header.Start, addressStart, addressStart)
		binFile.Header.Start = uint16(addressStart)
	}

	if len(tapHeaderName) == 0 {
		var array []byte = []byte(extractName(fileInName))
		if len(array) > 16 {
			array = array[:16]
		}
		for i, c := range array {
			if c < ' ' && c > '~' {
				array[i] = '.'
			}
		}
		tapHeaderName = string(array)
	}

	fmt.Printf(`
   Input file : %s
  Output file : %s
     TAP Name : %s
Start address : %d (&O%o)
   Turbo mode : %s
    Amplifier : %s
 `, fileInName, fileOutName, tapHeaderName, addressStart, addressStart, asOnOff(flagTurboMode), asOnOff(flagAmplify))

	fmt.Printf(`
   Values from BIN file header:
        Start  : %d (&O%o)
	Length : %d (&O%o)
`, binFile.Header.Start, binFile.Header.Start, binFile.Header.Length, binFile.Header.Length)

	checksum, err := bkbin.WriteWav(fileOutName, tapHeaderName, flagTurboMode, flagAmplify, &binFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("WAV file has been created successfully as '%s', checksum is #%X\n", fileOutName, checksum)
}
