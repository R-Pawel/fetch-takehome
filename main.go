package main

import (
	"github.com/R-Pawel/fetch-takehome/internal/routes"
)

func main() {
	router := routes.NewRouter()
	router.Run("0.0.0.0:8080")
}
