package modules

import (
	"fmt"
)

func init() {
	fmt.Println("Doing some init beep beep!")
}

func SwitchDevice(id string, command string) {
	fmt.Printf("test command", command)
}
