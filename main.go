package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"os"
)

// Create struct for settings
// TODO: move to some more usefull location
type Config struct {
	Name		string
	Port		string
	HTMLPath	string

}

func main() {

	fmt.Println("[SERVER STARTED]")

	// Load settings from conf.json
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	conf := Config{}

	// Handle parsing errors
	err := decoder.Decode(&conf)

	fmt.Println(conf);

	if err != nil {
	  fmt.Println("error:", err)
	}


	// Start testing go server
    http.Handle("/", http.FileServer(http.Dir(conf.HTMLPath)))
    http.ListenAndServe(conf.Port, nil)
}
