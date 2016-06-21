package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"encoding/json"

)

type Room struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	DeviceList []int  `json:"devices"`
}

// NewRoom creates a new room entry in the database
func NewRoom(t *Room) error {
	fmt.Println(t)
	return DB.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("rooms"))

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

// ListDevices returns a array of all device structs in the devices bucket
// Since amount should be 100 max for the beginning no need for splitting entries.
func ListRooms() []Room {
	var rooms []Room
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("rooms"))

		b.ForEach(func(k, v []byte) error {
			r := Room{}

			// Decode json value of cell to Device struct type
			err := json.Unmarshal(v, &r)
			if err != nil {
				return nil
			}
			rooms = append(rooms, r)
			fmt.Println(k, r)

			return nil
		})

		return nil
	})

	return rooms
}

// This part is currently not working
func GetRoom(ID int) Room {
	var room Room
	DB.View(func(tx *bolt.Tx) error {
	    b := tx.Bucket([]byte("rooms"))
		r := Room{}

	    v := b.Get(itob(ID))

		err := json.Unmarshal(v, &r)
		if err != nil {
			return nil
		}

	    room = r

	    return nil
	})

	return room
}

func deleteRoom(ID int) {


}
