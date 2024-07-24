package main

import (
	"github/oxiginedev/gostarter/cmd"
	"os"
)

func main() {
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
