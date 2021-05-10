package main

import (
	"github.com/matsuri-tech/golint-extra/rules"
	"golang.org/x/tools/go/analysis/multichecker"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if os.Getenv("DEBUG") != "true" {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stdout)
	}

	multichecker.Main(
		rules.NewZeroValueStruct(),
	)
}
