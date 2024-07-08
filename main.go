package main

import (
	"fmt"
	logging "goservice/logging"
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
App ID = ttnreceiver_maint
App Key = bwiawoitmisbstdwblhyh
*/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Server closing")
	}
	// Start up logging
	logging.Logging()

	server := NewServer(os.Getenv("HTTP_PORT"))

	fmt.Printf("Starting Server on port %v\n", server.Addr)
	log.Printf("Starting Server on port %v\n", server.Addr)

	go func() {
		// This starts the HTTP server
		err := server.ListenAndServe()

		if err != nil {

			log.Fatalln("server exiting:", err.Error())

		}
	}()
	//wait shutdown
	server.WaitShutdown()

	log.Println("server closing")

}
