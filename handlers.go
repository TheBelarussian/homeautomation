package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	//color "github.com/fatih/color"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// DeviceControlHandler handles incomming control commands from client
func DeviceControlHandler(w http.ResponseWriter, r *http.Request) {

	var command GenericCommand

	// Read data from client
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	// Check for error
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	// Parse data from client
	err = json.Unmarshal(body, &command)

	fmt.Println("Calling DeviceControlHandler")
	fmt.Println(command)

	DispatchDeviceCommand(&command)
}

// DeviceCreateHandler handles new device request.
func DeviceCreateHandler(w http.ResponseWriter, r *http.Request) {
	var device Device

	// Read data from client
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	// Check for error
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	// Parse data from client
	err = json.Unmarshal(body, &device)

	// Check if all requered fields are filled.
	if device.Name == "" || device.Type == "" {
		fmt.Print("Invalid data from client. Name and Type must be defined for a new device.")
		err = errors.New("Some requered field in device struct are empty.")
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		return
	}

	fmt.Println("Data fromn client:", device)

	// Send response to client if new item was created
	_, err = device.SaveToDB()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

}

// DeviceDeleteHandler handles deletion of devices.
func DeviceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var device Device

	// Read data from client
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	// Check for error
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	// Parse data from client
	err = json.Unmarshal(body, &device)

	fmt.Print(device)

	// Check if all requered fields are filled.
	if device.ID == 0 {
		fmt.Print("Invalid data from client. ID is requered when removing a device.")
		err = errors.New("Some requered field in device struct are empty.")
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		return
	}

	fmt.Println("Data fromn client:", device)

	err = device.RemoveFromDB()

	if err != nil {
		color.Red("Could not remove device")
	}
}

// DeviceListHandler returns a list of all devices as a JSON array.
func DeviceListHandler(w http.ResponseWriter, r *http.Request) {
	devices := ListDevices()

	// Try to send list of devices to client
	data, err := json.Marshal(struct {
		Data []Device `json:"data"`
	}{devices})
	// If this fails responde with error code 500 for internal server error
	if err != nil {

		if Debug {
			fmt.Println("Error: Could not connvert devices slice into JSON")
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Internal server error
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// Respose to Request
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write(data)

	//fmt.Println(string(data))
}

// DeviceGetHandler handles request for data on a device with given id, returns data as JSON.
func DeviceGetHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Device id:", mux.Vars(r)["deviceID"])

	// Get required id.
	id, err := strconv.Atoi(mux.Vars(r)["deviceID"])

	// Handle if id was not parsed.
	if err != nil {
		fmt.Println("Could not parse the given device ID")
		return
	}

	device, err := GetDevice(id)

	// Try to send list of devices to client
	data, err := json.Marshal(device)
	// If this fails responde with error code 500 for internal server error
	if err != nil {

		if Debug {
			fmt.Println("Error: Could not connvert devices slice into JSON")
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Internal server error
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// Respose to Request
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write(data)

	if Debug {
		fmt.Println(string(data))
	}
}

// RoomCreateHandler handles creation of a new room.
func RoomCreateHandler(w http.ResponseWriter, r *http.Request) {
	var room Room

	// Read data from client
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	// Check for error
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Parse data from client
	if err := json.Unmarshal(body, &room); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	fmt.Println("Data fromn client:", room)

	// Send response to client if new item was created
	t := NewRoom(&room)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

// RoomListHandler lists all rooms as JSON data.
func RoomListHandler(w http.ResponseWriter, r *http.Request) {
	rooms := ListRooms()

	// Try to send list of devices to client
	data, err := json.Marshal(rooms)
	// If this fails responde with error code 500 for internal server error
	if err != nil {

		if Debug {
			fmt.Println("Error: Could not connvert devices slice into JSON")
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Internal server error
		w.WriteHeader(500)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// Respose to Request
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write(data)

	if Debug {
		fmt.Println(string(data))
	}
}
