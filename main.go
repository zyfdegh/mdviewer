package main

import (
	"log"
	"os"

	"github.com/zyfdegh/mdviewer/cmd"
)

func main() {
	//call console
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("execute cmd error: %v", err)
		os.Exit(-1)
	}
}
