package main

import (
	"github.com/matsuri-tech/golint-extra/rules"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		rules.NewZeroValueStruct(),
	)
}
