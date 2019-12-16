package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gluk256/crypto/crutils"
)

const PageSize = 1503

var (
	encode         bool
	inFile         bool
	outFile        bool
	outputFileName string
	src            []byte
)

func help() {
	fmt.Println("hex v.0.14.9")
	fmt.Println("USAGE: hex flags src [dst]")
	fmt.Println("\t h help")
	fmt.Println("\t e encode")
	fmt.Println("\t d decode")
	fmt.Println("\t i input from file")
	fmt.Println("\t o output to file")
}

func readFile(name string) []byte {
	f, err := os.Open(name)
	if err != nil {
		fmt.Printf("Failed to read file [%s]\n", name)
		fmt.Printf("Error: [%s]\n", err.Error())
		os.Exit(0)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	b := make([]byte, 1024*1024)
	_, err = r.Read(b)
	if err != nil {
		fmt.Printf("Failed to read file [%s]\n", name)
		fmt.Printf("Error: [%s]\n", err.Error())
		os.Exit(0)
	}
	return b
}

func processFlags() {
	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], "h") || strings.Contains(os.Args[1], "?") {
			help()
			os.Exit(0)
		}
	}

	if len(os.Args) < 3 {
		fmt.Println("Error: not enough parameters")
		help()
		os.Exit(0)
	} else if len(os.Args) > 3 {
		outputFileName = os.Args[3]
		outFile = true
	}

	flags := os.Args[1]
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
		encode = false // decode flag overrides other settings
	}

	if inFile {
		src = readFile(os.Args[2])
	} else {
		src = []byte(os.Args[2])
	}
}

func main() {
	processFlags()

	if encode {
		if outFile {
			s := fmt.Sprintf("%x", src)
			saveResult([]byte(s))
		} else {
			processResult(src)
		}
	} else {
		res, err := crutils.HexDecode(src)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		} else {
			saveResult(res)
		}
	}
}

func saveResult(res []byte) {
	if len(outputFileName) == 0 {
		t := time.Now().UTC().UnixNano()
		rand.Seed(t)
		outputFileName = fmt.Sprintf("hex-%x", rand.Int()+int(t))
	}
	err := ioutil.WriteFile(outputFileName, res, 0666)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func processResult(data []byte) {
	if len(data) < PageSize {
		fmt.Printf("%x\n", data)
		return
	}

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "icanon").Run()

	total := len(data)/PageSize + 1
	var pg int
	var b []byte = make([]byte, 1)
	for b[0] != byte(27) && b[0] != byte(113) { // esc or 'q'
		pg = processCommand(b[0], pg, total)
		printPage(data, pg)
		os.Stdin.Read(b)
	}
}

func processCommand(c byte, pg int, total int) int {
	if c == 45 || c == 55 || c == 56 || c == 57 {
		pg--
		if pg < 0 {
			pg = 0
		}
	} else if c == 43 || c == 46 || c == 48 || c == 49 || c == 50 || c == 51 || c == 10 {
		pg++
		if pg >= total {
			pg = total - 1
		}
	}
	return pg
}

func printPage(data []byte, pg int) {
	fmt.Printf("PAGE %d\n\n", pg)
	beg := pg * PageSize
	end := beg + PageSize
	if end > len(data) {
		end = len(data)
	}
	fmt.Printf("%x\n\n", data[beg:end])
}
