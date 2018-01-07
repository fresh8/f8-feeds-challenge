package main

import (
	"fmt"
	"log"

	"github.com/Roverr/f8-feeds-challenge/importer/config"
)

// Very minimal error logging
func logError(customMsg string, err error) {
	fmt.Printf("ERROR || %s || %s\n", customMsg, err.Error())
}
func logInfo(customMsg string) {
	fmt.Printf("INFO || %s", customMsg)
}

func main() {
	// Initialize configuration and API
	_, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
}
