package main

import (
	application "Distributed_calculator/internal/app"
	"fmt"
)

func main() {
	app := application.New()
	err := app.RunServer()
	fmt.Println("Run Server")
	if err != nil {
		return
	}
}
