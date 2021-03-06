package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// Logger function for http requests
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// Route holds inforamation about which action to perform for a given route.
// TODO: make local
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes hold information for all route elements.
// Therefore a Routes is a collection of Route structs
type Routes []Route

// Create a router for the http.ListenAndServe.
// The routes for the router come from the router collection below.
func NewRouter() *mux.Router {

	if Debug {
		fmt.Println("[Creating Router]")
	}

	router := mux.NewRouter()

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		if Debug {
			fmt.Println("Endpoint: ", route.Name, route.Pattern)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Add static file routing
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(Conf.HTMLPath)))
	if Debug {
		fmt.Println("[Creating Router]")
	}
	if Debug {
		color.Green("[Router: OK]")
	}

	return router
}

var prefix = "/api"

var routes = Routes{
	// Device based routes
	Route{
		"CreateDevice",
		"CREATE",
		prefix + "/device/add",
		DeviceCreateHandler,
	},
	Route{
		"DeleteDevice",
		"DELETE",
		prefix + "/device/delete",
		DeviceDeleteHandler,
	},
	Route{
		"ListDevices",
		"GET",
		prefix + "/devices/list",
		DeviceListHandler,
	},

	Route{
		"GetDevice",
		"GET",
		prefix + "/device/{deviceID}",
		DeviceGetHandler,
	},

	Route{
		"ControlDevice",
		"POST",
		prefix + "/device/{deviceID}",
		DeviceControlHandler,
	},

	// Room based routes
	Route{
		"CreateRoom",
		"POST",
		prefix + "/rooms/add",
		RoomCreateHandler,
	},
	Route{
		"ListRooms",
		"GET",
		prefix + "/rooms/list",
		RoomListHandler,
	},
}
