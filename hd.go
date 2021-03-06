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
	"unicode/utf8"
)

const (
	OffLen  = 6       // Default length of the output Address/Offset field
	Version = "1.0.2" // Version Number
)

// Flag and other variables
var (
	inFile   string
	offBegin int
	offEnd   int
	lineLen  int
	innerLen int
	goDump   bool
	quiet    bool
	h        bool
	nobu     bool
	hexUpper bool
	edgeMark string
	inString string
	//
	argFname  string
	fileLen   = -1
	addrFlen  = -1
	minOffLen = -1
	version   bool
	//
	normLineLen = -1
	lastLineLen = -1
	//                      1         2         3         4
	//             1234567890123456789012345678901234567890
	pad    = "                                        "
	offFmt = "%016x"
	lbFmt  = "%s%02x"
	//
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

	flag.StringVar(&inString, "inString", "",
		"input string to dump.")

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

	flag.BoolVar(&nobu, "nobu", false,
		"if true, show bad UTF8 in character output.")

}

/*
Error checker.  Quit if an error is not nil.
*/
func checkError(e error, ds string) {
	if e != nil {
		fmt.Printf("\n%s %s\n\n", ds, e)
		if !quiet {
			fmt.Println("DumpFile Ends, RC:", 1)
		}
		os.Exit(1)
	}
}

/*
Set initial digit count of the address/offset field.
*/
func hexDigitCount(i int) {
	addrFlen = 1
	for {
		i = i / 16
		if i == 0 {
			return
		}
		addrFlen++
	}
	addrFlen++
	if addrFlen < OffLen {
		addrFlen = OffLen
	}
	if addrFlen < minOffLen {
		addrFlen = minOffLen
	}
}

/*
Set the size of the file being processed.
*/
func setFileLen(f *os.File) {
	fi, err := f.Stat()
	checkError(err, "Stat Error ==>")
	fileLen = int(fi.Size())
	hexDigitCount(fileLen)
	// fmt.Printf("Hex Digit Count: %d\n", addrFlen)
}

/*
Open a user specified file to get an io.Reader.
*/
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

/*
High level return of an io.Reader.
*/
func getReader() io.Reader {
	fa := flag.Args()
	if len(fa) >= 1 {
		argFname = fa[0]
	}
	if inString != "" {
		hexDigitCount(len(inString))
		return strings.NewReader(inString)
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

/*
Use format provided in the encoding/hex package.
*/
func goFormatDump(r io.Reader) {
	// Dump
	buff, err := ioutil.ReadAll(r)
	checkError(err, "ReadAll error ==>")
	fmt.Printf("%s", hex.Dump(buff))
	//
	return
}

/*
Print the address/offset field.
*/
func printOffset(o int) {
	had := fmt.Sprintf(offFmt, o)
	fmt.Printf("%s  ", had[16-addrFlen:])
}

/*
Get a buffer of blanks.
*/
func blankBuf(l int) []byte {
	s := strings.Repeat(" ", l)
	return []byte(s)
}

/*
Print the left buffer (left side of output line).
*/
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
		//
		/*
			if len(pad) < normLineLen-lastLineLen {
				fmt.Println()
				fmt.Println("lp", len(pad), "nx", normLineLen, lastLineLen)
			}
		*/
		fmt.Print(pad[:normLineLen-lastLineLen])
	}
	// fmt.Println("lll", lastLineLen, "nll", normLineLen)
	fmt.Print(edgeMark)
}

/*
Print the right buffer (right side of output line).
*/
func printRightBuffer(br int, ib []byte) {
	bb := blankBuf(br)
	for i := 0; i < br; i++ {
		bb[i] = ib[i]
		if bb[i] < byte(0x20) {
			bb[i] = byte('.')
		}
		if bb[i] == byte(0x7f) {
			bb[i] = byte('.')
		}
	}
	rhs := string(bb)
	// fmt.Printf("\nb16: %#x\n", rhs[0:16])
	if nobu {
		if !utf8.ValidString(rhs) {
			rr := []rune(rhs)
			wrs := make([]rune, 0)
			for _, v := range rr {
				if utf8.ValidRune(v) {
					wrs = append(wrs, v)
				} else {
					wrs = append(wrs, ' ')
				}
				rhs = string(wrs)
			}
		}
		//
	}
	fmt.Print(rhs)
	fmt.Print(edgeMark)
	//
}

/*
Dump file contents in hex format.
*/
func main() {
	flag.Parse() // Parse all flags

	// fmt.Println("nobu: ", nobu)
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
	for {
		readLen := lineLen
		if offEnd > 0 && roff+readLen > offEnd {
			readLen = offEnd - roff + 1
		}
		ib := blankBuf(readLen)
		//fmt.Println("ReadLen is now:", readLen)
		br, _ := io.ReadAtLeast(r, ib, readLen)
		//fmt.Println("Actual Read Length:", br)
		if br == 0 {
			break
		}
		//fmt.Println("Start print offset")
		printOffset(roff)
		printLeftBuffer(br, ib)
		printRightBuffer(br, ib)
		roff += lineLen
		if offEnd > 0 && roff > offEnd {
			fmt.Println()
			break
		}
		fmt.Println()
	}
	if !quiet {
		fmt.Println("DumpFile Ends....")
	}
}
