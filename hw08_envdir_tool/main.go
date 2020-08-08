package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Wrong arguments count: min 2 args")
		os.Exit(111)
	}
	dir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(111)
	}
	os.Exit(RunCmd(command, env))
}
