package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gluk256/crypto/crutils"
)

var (
	encode bool
	inFile bool
	outFile bool
	outputFileName string
)

func help() {
	fmt.Println("hex v.0.13")
	fmt.Println("USAGE: hex [flags] [src]")
	fmt.Println("\t h help")
	fmt.Println("\t e encode")
	fmt.Println("\t d decode")
	fmt.Println("\t i input from file")
	fmt.Println("\t o output to file")
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
		if strings.Contains(os.Args[1], "h") {
			help()
			return
		}
	}

	if len(os.Args) < 3 {
		fmt.Println("Error: not enough parameters")
		help()
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
		var err error
		res, err = crutils.HexDecode(src)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
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
