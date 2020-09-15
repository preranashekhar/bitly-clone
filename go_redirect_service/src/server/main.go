package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}

	serverInit(":" + port)
}