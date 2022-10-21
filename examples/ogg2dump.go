package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/drgolem/go-ogg/ogg"
)

func main() {
	fmt.Println("example decode ogg file")

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: ogg2dump <infile.ogg>")
		return
	}

	inFile := os.Args[1]
	fmt.Printf("infile: %s\n", inFile)

	f, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	oggReader, err := ogg.NewOggReader(reader)
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return
	}
	defer oggReader.Close()

	pktCnt := 0
	for oggReader.Next() {
		p, err := oggReader.Scan()
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
		pktCnt++
		fmt.Printf("packet [%d] len: %d\n%s\n", pktCnt, len(p), hex.Dump(p))
	}
}
