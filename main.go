package main

import (
	"io"
	"log"
	"os"

	"github.com/matsuri-tech/golint-extra/rules"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	if os.Getenv("DEBUG") != "true" {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stdout)
	}

	multichecker.Main(
		rules.NewZeroValueStruct(),
	)
}
