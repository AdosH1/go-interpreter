package main

import (
	"fmt"
	"interpreter/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s, welcome to my Go Interpreter!\n", user.Username)
	fmt.Printf("Try some commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
