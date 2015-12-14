package main

import (
	"fmt"
	"encoding/json"
	"os"
	"github.com/go-martini/martini"
)

// Create struct for settings
// TODO: move to some more usefull location
type Config struct {
	Name		string
	Port		string
	HTMLPath	string

}

func main() {
	m := martini.Classic()
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

	m.Get("/test", func() string {
	  return "Hello world!"
	})
	m.Use(martini.Static("homeauto-client/dist"))
	m.Run()
}
