package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
)

// Device is a struct representation of a device.
type Device struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	RoomID int    `json:"roomid"`
}

// NewDevice creates a new device entry in the database.
// Return an error and the id of the created device.
func NewDevice(t *Device) (int, error) {
	fmt.Println(t)
	return t.ID, DB.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("devices"))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		t.ID = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(t.ID), buf)
	})
}

// RemoveDevice removes a device from the database.
func RemoveDevice(id int) error {
	if Debug {
		color.Red("Removing device ID:", id)
	}

	// Start a bolt transaction.
	return DB.Update(func(tx *bolt.Tx) error {
		// Try to delete device with the given id from the bucket.
		v := tx.Bucket([]byte("devices")).Delete(itob(id))
		fmt.Print("response", v)
		return v
	})
}

// ListDevices returns a array of all device structs in the devices bucket
// Since amount should be 100 max for the beginning no need for splitting etries.
func ListDevices() []Device {
	var devices []Device
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))

		b.ForEach(func(key, v []byte) error {
			device := Device{}
			// Decode json value of cell to Device struct type
			err := json.Unmarshal(v, &device)
			if err != nil {
				return nil
			}
			devices = append(devices, device)
			fmt.Println(key, device)

			return nil
		})

		return nil
	})
	return devices
}

// GetDevice return a device. If device is not found or an error occurs while getting the device
// an error is beeing thrown.
func GetDevice(id int) (Device, error) {
	var device Device
	err := DB.View(func(tx *bolt.Tx) error {
		// Try to get the device from the bucket.
		value := tx.Bucket([]byte("devices")).Get(itob(id))

		// If no device found return error.
		if value == nil {
			return errors.New("No device found with given id.")
		}

		// Parse value into the device struct.
		return json.Unmarshal(value, &device)
	})
	return device, err
}
