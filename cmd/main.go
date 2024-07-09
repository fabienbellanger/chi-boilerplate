package main

import (
	"chi_boilerplate/pkg/infrastructure/cli"
	"log"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatalln(err)
	}
}
