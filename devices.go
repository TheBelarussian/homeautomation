package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

// Device is a struct representation of a device.
type Device struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	RoomID int    `json:"roomid"`
}

// NewDevice creates a new device entry in the database
func NewDevice(t *Device) error {
	fmt.Println(t)
	return DB.Update(func(tx *bolt.Tx) error {
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

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// ListDevices returns a array of all device structs in the devices bucket
// Since amount should be 100 max for the beginning no need for splitting etries.
func ListDevices() []Device {
	var devices []Device
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))

		b.ForEach(func(k, v []byte) error {
			d := Device{}
			// Decode json value of cell to Device struct type
			err := json.Unmarshal(v, &d)
			if err != nil {
				return nil
			}
			devices = append(devices, d)
			fmt.Println(k, d)

			return nil
		})

		return nil
	})
	return devices
}

//
func GetDevice(ID int) Device {
	device := Device{}
	
	// Get device from DB
	return device
}
