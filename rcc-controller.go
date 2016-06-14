package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/stianeikeland/go-rpio"
)

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	blinkPin = rpio.Pin(4)
	rccPin   = rpio.Pin(3)
	//pulseLength = 300
)

func testRCCSend() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	rccPin.Output()

	// Toggle pin 20 times
	for x := 0; x < 20; x++ {
		state := x % 2
		color.Green("Sending: ", state)
		if state == 1 {
			rccPin.High()
		} else {
			rccPin.Low()
		}
		time.Sleep(300)
	}
}

func testGPIO() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	blinkPin.Output()

	// Toggle pin 20 times
	for x := 0; x < 20; x++ {
		color.Green("Blinking...")
		blinkPin.Toggle()
		time.Sleep(time.Second / 5)
	}
}
