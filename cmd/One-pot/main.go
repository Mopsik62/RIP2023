package main

import (
	"awesomeProject1/internal/api"
	"log"
)

func main() {
	log.Println("Application start!")

	a := app.New()
	a.StartServer()

	log.Println("Application terminated!")
}
