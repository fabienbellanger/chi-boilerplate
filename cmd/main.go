package main

import (
	"fmt"

	"test_chi/pkg/infrastructure/chi_router"
)

func main() {
	fmt.Println("Test HTTP server")

	chi_router.Start()
}
