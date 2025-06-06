package main

import (
	"fmt"
	"os"

	"github.com/IBM/openkommander/pkg/cli"
)

func main() {
	var rootCmd = cli.Init()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
