package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
)

var (
	encode bool
	inFile bool
	outFile bool
	outputFileName string
)

func help() {
	fmt.Println("hex v.0.12")
	fmt.Println("USAGE: hex [flags] [src]")
	fmt.Println("\t h help")
	fmt.Println("\t e encode")
	fmt.Println("\t d decode")

}

func readFile(name string) []byte {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Printf("Error: can not read file [%s]\n", name)
		os.Exit(0)
	}
	return b
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1][0] == 'h' {
			help()
			return
		}
	}

	if len(os.Args) < 3 {
		fmt.Println("Error: not enough parameters")
		return
	} else if len(os.Args) > 3 {
		outputFileName = os.Args[3]
	}

	var src []byte
	flags := os.Args[1]
	s := os.Args[2]
	if strings.Contains(flags, "h") {
		help()
		return
	}
	if strings.Contains(flags, "?") {
		help()
		return
	}
	if strings.Contains(flags, "e") {
		encode = true
	}
	if strings.Contains(flags, "i") {
		inFile = true
	}
	if strings.Contains(flags, "o") {
		outFile = true
	}
	if strings.Contains(flags, "d") {
		// decode overrides other settings
		encode = false
	}

	if inFile {
		src = readFile(s)
	} else {
		src = []byte(s)
	}

	var res []byte
	if encode {
		res = hexEncode(src)
	} else {
		res = hexDecode(src)
	}

	if outFile {
		outputResult(res)
	} else {
		fmt.Printf("%s\n", string(res))
	}
}

func outputResult(res []byte) {
	if len(outputFileName) == 0 {
		t := time.Now().UTC().UnixNano()
		rand.Seed(t)
		outputFileName = fmt.Sprintf("hex-%x", rand.Int() + int(t))
	}
	err := ioutil.WriteFile(outputFileName, res, 0666)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func hexEncode(src []byte) []byte {
	s := fmt.Sprintf("%x", src)
	return []byte(s)
}

func convert(b byte) int {
	if b >= 48 && b <= 57 {
		return int(b - 48)
	}
	if b >= 65 && b <= 70 {
		return int(b - 65) + 10
	}
	if b >= 97 && b <= 102 {
		return int(b - 97) + 10
	}
	return -1
}

func hexDecode(src []byte) []byte {
	for i := len(src) - 1; i >=0; i-- {
		if src[i] > 32 && src[i] < 128 {
			break
		} else {
			src = src[:len(src) - 1]
		}
	}

	sz := len(src)
	if sz % 2 == 1 {
		fmt.Printf("Error decoding: odd src size %d\n", sz)
		os.Exit(0)
	}

	var dst []byte
	for i := 0; i < sz; i += 2 {
		a := convert(src[i])
		b := convert(src[i+1])
		if a < 0 || b < 0 {
			fmt.Printf("Error deocding: illegal byte %s\n", string(src[i:i+2]))
			os.Exit(0)
		}
		dst = append(dst, byte(16*a+b))
	}
	return dst
}


