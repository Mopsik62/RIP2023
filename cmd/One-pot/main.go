package main

import (
	"awesomeProject1/internal/pkg/app"
	"context"
	"log"
)

// @title One-pot
// @version 0.0-0
// @description One-pot synthesis

// @host 0.0.0.0:8000
// @schemes http
// @BasePath /

func main() {
	log.Println("Application start!")

	a, err := app.New(context.Background())
	if err != nil {
		log.Println(err)

		return
	}

	a.StartServer()

	log.Println("Application terminated!")
}
