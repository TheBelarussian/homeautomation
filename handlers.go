package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	//color "github.com/fatih/color"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)

	fmt.Println(session.Values[42])
}

// Handle new device request
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
	_, err = NewDevice(&device)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

}

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

	err = RemoveDevice(device.ID)

	if err != nil {
		color.Red("Could not remove device")
	}
}

func DeviceListHandler(w http.ResponseWriter, r *http.Request) {
	devices := ListDevices()

	// Try to send list of devices to client
	data, err := json.Marshal(struct {
		Data []Device `json: "data"`
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

	fmt.Println(string(data))
}

func DeviceGetHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Device id:", mux.Vars(r)["deviceID"])
	device, err := GetDevice(0)

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

	fmt.Println(string(data))
}

/* ADD SOME COMMENTS*/
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

// Lists all rooms
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

	fmt.Println(string(data))
}
