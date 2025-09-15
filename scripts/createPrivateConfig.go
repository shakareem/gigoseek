package main

import (
	"log"
	"os"
)

func main() {
	public, err := os.Open("./configs/public.json")
	if err != nil {
		log.Fatalf("Failed to open public config file: %v", err)
	}
	defer public.Close()

	private, err := os.OpenFile("./configs/private.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open or create private config file: %v", err)
	}
	defer private.Close()

	_, err = private.ReadFrom(public)
	if err != nil {
		log.Fatalf("Failed to copy content from public to private config file: %v", err)
	}

	log.Println("Public config copied to private config successfully.")
}
