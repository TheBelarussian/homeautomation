package main

import (
	"fmt"
	"encoding/json"
	"os"
	"log"
	"os/exec"
	"github.com/go-martini/martini"
)

// Create struct for settings
// TODO: move to some more usefull location
type Config struct {
	Name			string `json:"name"`
	Port			string `json:"port"`
	HTMLPath		string `json:"wwwPath"`
	Raspremote 		string `json:"raspremote"`
	Test			string `json:"test"`
}

func main() {
	m := martini.Classic()
	fmt.Println("[READING CONFIG]")

	var settings Config;

	// then config file settings
	configFile, err := os.Open("conf.json")
	if err != nil {
	    log.Fatal(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&settings); err != nil {
	    log.Fatal(err)
	}

	fmt.Println(settings);

	if err != nil {
	  fmt.Println("error:", err)
	}


	// Start testing go server

	m.Get("/test/:device/:command", func(params martini.Params) string {
		fmt.Println(settings.Raspremote + " " + params["device"] + " " + params["command"]);
		testString	:= []string{"element1", "element2", "element3"}

		// Call send comand for 433 Mhz receiver
		cmd :=  settings.Raspremote + " " + params["device"] + " " + params["command"]

		exec.Command("sh","-c", cmd).Output()

    	json, _ := json.Marshal(testString)

		return string(json);
	})
	m.Use(martini.Static("homeauto-client/dist"))
	m.Run()
}
