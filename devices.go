package main

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

// Device is a struct representation of a device.
type Device struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	RoomID int    `json:"roomid"`
}

// GenericCommand proviedes a generalized structure for the command calls.
type GenericCommand struct {
	Type     string `json:"type"`
	TargetID int    `json:"id"`
	Command  string `json:"command"`
}

// GetID implements StorableElement interface.
func (d Device) GetID() int {
	return d.ID
}

// SetID implements StorableElement interface.
func (d Device) SetID(id int) StorableElement {
	d.ID = id
	return d
}

// ConvertToJSON implements StorableElement interface for converting Struct to JSON.
func (d Device) ConvertToJSON(id int) ([]byte, error) {
	if id != 0 {
		d.ID = id
	}
	return json.Marshal(d)
}

// RemoveFromDB implements StorableElement interface for removing from db.
func (d Device) RemoveFromDB() error {
	return RemoveFromDB("devices", d.ID)
}

// SaveToDB implements StorableElement interface for storing element in db.
func (d Device) SaveToDB() (int, error) {
	return StoreToDB("devices", d)
}

// DispatchDeviceCommand handles mapping of command to a given device type.
func DispatchDeviceCommand(command *GenericCommand) {
	// TODO: validate if user has rights to control device here...

	// Here we have to math a givent device Type to a control Type. Hence a device can be controlled by a number of ways.
	// In order to allow modularity for different types of devices we have to map a control function to each type.
	callback := GetDeviceTypeOrNil("RCCSimple")
	callback(command)
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
			//fmt.Println(key, device)

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
