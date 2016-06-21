package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"wladmeixner.de/homeautomat/modules"

	"github.com/boltdb/bolt"
	color "github.com/fatih/color"
	"github.com/gorilla/sessions"
)

// Create struct for settings
// TODO: move to some more usefull location
type Config struct {
	Name     string
	Port     string
	HTMLPath string
	Database string
}

// Global vars
var (
	Debug      bool
	DummyMode  bool
	ConfigPath string
	Conf       = Config{}
	// Construct better structure for DB (for testing purposes)
	DB *bolt.DB
	// Local vars
	// Test for cookie session management (not ready)
	store = sessions.NewCookieStore([]byte("something-very-secret"))
)

// Read in and parse the congig
// path: Path to config file
func initConf(path string) error {
	// Load settings from conf.json
	file, err := os.Open(path)
	decoder := json.NewDecoder(file)

	// Could not load config file.
	if err != nil {
		color.Red("ERROR:", err)

		if path != "conf.json" {
			fmt.Println("Looking in current directory")
			initConf("./conf.json")
		}
		return err
	}

	// Handle parsing errors
	err = decoder.Decode(&Conf)

	if err != nil {
		fmt.Println("error:", err.Error())
		return err
	}

	if Debug {
		color.Green("[SETTINGS: OK]")
		fmt.Println("Settings Data:", Conf)
	}

	return nil
}

// Init basic strctures in the database
func initDatabase(dbPath string) (*bolt.DB, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	color.Green("[DB: OK]")
	// Init buckets if not present
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("devices"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Could not create bucket in DB", err)
	}
	return db, nil
}

// Prepare programm for equit.
// Close server, write data and close sockets.
func controlledQuit(c chan os.Signal) {
	<-c
	quit()
}

func quit() {
	color.Green("Cleany closy worky worky...")
	// Close database.
	err := DB.Close()
	if err != nil {
		color.Red("Could not close DB ", err)
	}
	color.Yellow("Quiting programm")

	// Exit programm
	os.Exit(0)
}

// Reads in all command line arguments and
func main() {

	// Assign flag
	flag.BoolVar(&Debug, "debug", false, "Almost everything will be logged if this is true")

	// Assign flag
	flag.BoolVar(&DummyMode, "dummy", false, "Starts server in debug mode with dummy 433 MHz controller")

	// Assign flag
	flag.StringVar(&ConfigPath, "path", "conf.json", "Give path to the config pathr")

	// Parse flags
	flag.Parse()

	// TEST FLAG values
	if Debug {
		color.Yellow("[SERVER RUNNING IN DEBUG MODE!!]")
	}

	if DummyMode {
		color.Red("[SERVER RUNNING IN DUMMY MODE!]")
	}

	// Read config file
	err := initConf(ConfigPath)

	if err != nil {
		color.Red("ERROR while parsing config file.", err.Error())
		return
	}

	// Setup server quit action (^C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go controlledQuit(c)
	// Init bolt database
	DB, err = initDatabase(Conf.Database)
	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	if Debug {
		d := Device{0, "Test device", "Lamp", 22}

		deviceID, err := NewDevice(&d)

		if err != nil {
			color.Red("Failed to add new device to DB. This is kind of fatal so go panic, try not to cry, cry a lot!")
			quit()
		}

		color.Green("Device added with id:", deviceID)
		ListDevices()
		err = RemoveDevice(deviceID)

		if err != nil {
			color.Red("Failed to remove device from DB. This is kind of fatal so go panic and cry!")
			quit()
		}
		// Test GPIO

	}

	if !DummyMode {
		color.Yellow("[TESTING GPIO!!]")
		testGPIO()
		testRCCSend()
	}

	// Create Router (multiplexer)
	router := NewRouter()
	color.Green("[SERVER RUNNING]")
	log.Fatal(http.ListenAndServe(":"+Conf.Port, router))

}
