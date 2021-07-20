package main

import (
	"log"

	"github.com/mikydna/sports/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		log.Fatal(err)
	}
}
