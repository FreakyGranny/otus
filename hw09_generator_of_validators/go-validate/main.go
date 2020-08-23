package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("source argument not set")
		os.Exit(1)
	}
	sPath := os.Args[1]
	src, err := os.OpenFile(sPath, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Printf("FILE ERRROR %s\n", err)
		os.Exit(1)
	}
	defer src.Close()

	source, err := parseSource(src)
	if err != nil {
		fmt.Printf("PARSING ERRROR %s\n", err)
		os.Exit(2)
	}

	preparedData, err := prepareData(source)
	if err != nil {
		fmt.Printf("PREPARING ERRROR %s\n", err)
		os.Exit(3)
	}
	result, err := buildFile(preparedData)
	if err != nil {
		fmt.Printf("BUILDING ERRROR %s\n", err)
		os.Exit(4)
	}
	dst, err := os.Create(strings.Replace(sPath, ".go", "_validation_generated.go", 1))
	if err != nil {
		fmt.Printf("CREATING FILE ERRROR %s\n", err)
		os.Exit(5)
	}
	defer dst.Close()
	_, err = dst.Write(result)
	if err != nil {
		fmt.Printf("WRITING ERRROR %s\n", err)
		os.Exit(6)
	}
}
