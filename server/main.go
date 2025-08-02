package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

}