package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

// StorableElement must be implemented to use most of db functions and be called outside of it.
type StorableElement interface {
	SaveToDB() (int, error)
	RemoveFromDB() error
	SetID(int) StorableElement
	GetID() int
	ConvertToJSON(int) ([]byte, error)
}

// StoreToDB adds given data to the given bucket. The ID is choosen sequentially.
func StoreToDB(bucket string, data StorableElement) (int, error) {
	return data.GetID(), DB.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte(bucket))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()

		// Marshal user data into bytes.
		buf, err := data.ConvertToJSON(int(id))
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(int(id)), buf)
	})
}

// RemoveFromDB removes key,value in a given bucket if existing.
func RemoveFromDB(bucket string, id int) error {
	return DB.Update(func(tx *bolt.Tx) error {
		// Try to delete device with the given id from the bucket.
		v := tx.Bucket([]byte(bucket)).Delete(itob(id))
		fmt.Print("response", v)
		return v
	})
}

// GetElementData return byte value for the given bucket, key pair if existing.
func GetElementData(bucket string, id int) ([]byte, error) {
	var data []byte

	err := DB.View(func(tx *bolt.Tx) error {
		// Try to get the device from the bucket.
		data := tx.Bucket([]byte(bucket)).Get(itob(id))

		// If no device found return error.
		if data == nil {
			return errors.New("No device found with given id.")
		}

		// Return nil if device was valid.
		return nil
	})

	// Return data and error to caller.
	return data, err
}

// GetDataToStruct gets value for given bucket and get and parses it into a given struct.
// The Struct must at least implement StorableElement interface for this to work.
func GetDataToStruct(bucket string, id int, element *StorableElement) error {
	// Get value as byte array.
	data, err := GetElementData(bucket, id)

	// Check if error was returned.
	if err != nil {
		return err
	}

	// Parse data to given implementation of interface.
	return json.Unmarshal(data, element)
}