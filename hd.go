/*
	Dump a file in hex.
*/
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	OffLen  = 6       // Default length of the output Address/Offset field
	Version = "1.0.2" // Version Number
)

// Flag variables
var (
	inFile   string
	offBegin int
	offEnd   int
	lineLen  int
	innerLen int
	goDump   bool
	quiet    bool
	h        bool
	hexUpper bool
	edgeMark string
	//
	argFname  string
	fileLen   = -1
	addrFlen  = -1
	minOffLen = -1
	version   bool
	//
	normLineLen = -1
	lastLineLen = -1
	pad         = "                           "
	offFmt      = "%016x"
	lbFmt       = "%s%02x"
)

// Main initialization, set flags up
func init() {
	flag.StringVar(&edgeMark, "edgeMark", "|",
		"single character at edges of the right hand character block.")

	flag.BoolVar(&goDump, "goDump", false,
		"if true, use standard go encoding/hex/Dump.")

	flag.BoolVar(&h, "h", false, "print usage message.")

	flag.BoolVar(&hexUpper, "hexUpper", false,
		"if true, print upper case hex.")

	flag.StringVar(&inFile, "inFile", "",
		"input file name.  Argument 0 may also be used.")

	flag.IntVar(&innerLen, "innerLen", 4,
		"dump line inner area byte count.")

	flag.IntVar(&lineLen, "lineLen", 16,
		"dump line total byte count.")

	flag.IntVar(&minOffLen, "minOffLen", -1,
		"minimum lenght of the offset field.")

	flag.IntVar(&offBegin, "offBegin", 0,
		"begin dump at file offset.")

	flag.IntVar(&offEnd, "offEnd", -1,
		"end dump at file offset.")

	flag.BoolVar(&quiet, "quiet", false,
		"if true, suppress header/trailer/informational messages.")

	flag.BoolVar(&version, "version", false,
		"if true, display program version and exit.")
}

func checkError(e error, ds string) {
	if e != nil {
		fmt.Printf("\n%s %s\n\n", ds, e)
		if !quiet {
			fmt.Println("DumpFile Ends, RC:", 1)
		}
		os.Exit(1)
	}
}

func hexDigitCount(i int) {
	addrFlen = 1
	for {
		i = i / 16
		if i == 0 {
			return
		}
		addrFlen++
	}
}

func setFileLen(f *os.File) {
	fi, err := f.Stat()
	checkError(err, "Stat Error ==>")
	fileLen = int(fi.Size())
	hexDigitCount(fileLen)
	addrFlen++
	if addrFlen < OffLen {
		addrFlen = OffLen
	}
	if addrFlen < minOffLen {
		addrFlen = minOffLen
	}
	// fmt.Printf("Hex Digit Count: %d\n", addrFlen)
}

func fileInit(fn, ed string) io.Reader {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0644)
	checkError(err, ed+" Open Error ==>")
	setFileLen(f)
	if offBegin > 0 {
		_, err := f.Seek(int64(offBegin), io.SeekStart)
		checkError(err, "Seek Error ==>")
	}
	return f
}

func getReader() io.Reader {
	fa := flag.Args()
	if len(fa) >= 1 {
		argFname = fa[0]
	}
	if inFile == "" && argFname == "" {
		addrFlen = OffLen  // Arbitrary, file size is unknown
		if minOffLen > 0 { // Let user specify length
			if addrFlen < minOffLen {
				addrFlen = minOffLen
			}
		}
		return os.Stdin
	}
	if inFile != "" {
		return fileInit(inFile, "inFile")
	}
	return fileInit(argFname, "argFname")
}

func goFormatDump(r io.Reader) {
	// Dump
	buff, err := ioutil.ReadAll(r)
	checkError(err, "ReadAll error ==>")
	fmt.Printf("%s", hex.Dump(buff))
	//
	return
}

func printOffset(o int) {
	had := fmt.Sprintf(offFmt, o)
	fmt.Printf("%s  ", had[16-addrFlen:])
}

func blankBuf(l int) []byte {
	s := strings.Repeat(" ", l)
	return []byte(s)
}

func printLeftBuffer(br int, ib []byte) {
	nol := (lineLen / innerLen) + 1
	os := ""
	noff := 0
leftFor:
	for no := 0; no < nol; no++ {
		for ni := 0; ni < innerLen; ni++ {
			if noff < br {
				nbi := int(ib[noff])
				os = fmt.Sprintf(lbFmt, os, nbi)
			} else {
				os = os + "  " // Add two blanks here
			}
			if normLineLen < len(os) {
				normLineLen = len(os)
			}
			noff++
			if noff > br {
				// os = os + "          " // Alignment need fixed.
				os = os + " " // End of left (center) data
				lastLineLen = len(os) - 1
				break leftFor
			}
		}
		os = os + " " // Blank at end of inner
	}
	fmt.Print(os)
	if lastLineLen < normLineLen {
		fmt.Print(pad[:normLineLen-lastLineLen])
	}
	// fmt.Println("lll", lastLineLen, "nll", normLineLen)
	fmt.Print(edgeMark)
}

func printRightBuffer(br int, ib []byte) {
	bb := blankBuf(br)
	for i := 0; i < br; i++ {
		bb[i] = ib[i]
		if bb[i] < byte(0x20) {
			bb[i] = byte('.')
		}
	}
	fmt.Print(string(bb))
	fmt.Print(edgeMark)
}

func main() {
	flag.Parse() // Parse all flags

	if h {
		flag.PrintDefaults()
		return
	}
	if version {
		fmt.Printf("hd Version: %s\n", Version)
		return
	}

	if !quiet {
		fmt.Println("DumpFile Starts....")
	}
	// fmt.Println("Line Length:", lineLen)
	r := getReader()
	if goDump {
		goFormatDump(r)
		if !quiet {
			fmt.Println("DumpFile Ends....")
		}
		return
	}
	if offEnd > 0 && offEnd <= offBegin {
		fmt.Printf("Offset Error: offEnd(%d) must be > offBegin(%d)\n",
			offEnd, offBegin)
		if !quiet {
			fmt.Println("DumpFile Ends, RC:", 2)
		}
		os.Exit(2)
	}
	if hexUpper {
		offFmt = "%016X"
		lbFmt = "%s%02X"
	}
	//
	roff := offBegin
	needcr := false
	for {
		readLen := lineLen
		if offEnd > 0 && roff+readLen > offEnd {
			readLen = offEnd - roff + 1
		}
		ib := blankBuf(readLen)
		// fmt.Println("ReadLen is now:", readLen)
		br, _ := io.ReadAtLeast(r, ib, readLen)
		// fmt.Println("Actual Read Length:", br)
		if br == 0 {
			break
		}
		printOffset(roff)
		printLeftBuffer(br, ib)
		printRightBuffer(br, ib)
		roff += lineLen
		if offEnd > 0 && roff > offEnd {
			needcr = true
			break
		}
		fmt.Println()
	}
	if needcr {
		fmt.Println()
	}
	if !quiet {
		fmt.Println("DumpFile Ends....")
	}
}
