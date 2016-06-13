package main

type Room struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	DeviceList []int  `json:"devices"`
}

var rooms []Room

var currentRoomId int = 0

func NewRoom(t Room) Room {
	currentRoomId += 1
	t.ID = currentRoomId
	rooms = append(rooms, t)
	return t
}

func ListRooms() []Room {
	return rooms
}
