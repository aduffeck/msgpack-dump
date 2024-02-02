package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/shamaton/msgpack/v2"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: [options] msgpack-dump <inputfile>\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var format string
	flag.StringVar(&format, "format", "plain", `output format ("plain" or "json")`)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 || len(args) > 1 {
		usage()
	}

	bytes, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("failed to read input file")
		os.Exit(1)
	}

	if len(bytes) == 0 {
		fmt.Println("input file is empty")
		os.Exit(0)
	}

	if (bytes[0] >= 0x80 && bytes[0] <= 0x8f) /* fixmap */ || bytes[0] == 0xde /* map 16 */ || bytes[0] == 0xdf /* map 32 */ {
		dumpMap(bytes, format)
		return
	}
	if (bytes[0] >= 0x90 && bytes[0] <= 0x9f) /* fixarray */ || bytes[0] == 0xdc /* array 16 */ || bytes[0] == 0xdd /* array 32 */ {
		dumpArray(bytes, format)
		return
	}

	fmt.Printf("Only map and array types are supported at that point. Type: %d\n", bytes[0])
	os.Exit(1)
}

func dumpMap(bytes []byte, format string) {
	m := map[interface{}]interface{}{}
	err := msgpack.Unmarshal(bytes, &m)
	if err != nil {
		fmt.Printf("Error unmarshalling the messagepack data: %s\n", err)
		os.Exit(1)
	}
	if format == "json" {
		jsonMap := map[string]string{}
		for k, v := range m {
			jsonMap[fmt.Sprintf("%s", k)] = string(fmt.Sprintf("%s", v))
		}
		b, err := json.MarshalIndent(jsonMap, "", "  ")
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		for k, v := range m {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
}

func dumpArray(bytes []byte, format string) {
	a := []interface{}{}
	err := msgpack.Unmarshal(bytes, &a)
	if err != nil {
		fmt.Printf("Error unmarshalling the messagepack data: %s\n", err)
		os.Exit(1)
	}

	if format == "json" {
		b, err := json.MarshalIndent(a, "", "  ")
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		for _, v := range a {
			fmt.Printf("%s\n", v)
		}
	}
}

func safeString(in interface{}) string {
	s := string(in)
	if _, ok := s.([]byte) {

	}
	return s
}
