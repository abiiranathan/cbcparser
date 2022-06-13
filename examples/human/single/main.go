package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abiiranathan/cbcparser/cbcparser"
	"github.com/abiiranathan/cbcparser/cbcparser/human"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <cbc_file> <normal_ranges_file>\n", os.Args[0])
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("open error: %s\n", err)
	}

	f2, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatalf("open error: %s\n", err)
	}

	defer f.Close()
	defer f2.Close()

	parser := human.New()

	// Parse with normal ranges
	normal_ranges, err := cbcparser.ReadNormalRanges(f2)
	if err != nil {
		log.Fatalf("read normal ranges error: %s\n", err)
	}

	w, err := parser.Parse(f, normal_ranges)

	if err != nil {
		log.Fatalf("parse error: %s\n", err)
	}

	err = w.Write(os.Stdout, cbcparser.JSONIndent)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
