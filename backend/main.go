package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	version = "1.0.1"
	sysLog  *log.Logger
	errLog  *log.Logger
)

// apiResponse is used to send uniform response structure.
type apiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

// sendEnvelope is used to send success response based on format defined in apiResponse
func sendEnvelope(w http.ResponseWriter, code int, message string, data interface{}) {
	// Standard marshalled envelope for success.
	a := apiResponse{
		Status:  "success",
		Data:    data,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		errLog.Panicf("Quitting %s", err)
	}
}

func initLogger() {
	sysLog = log.New(os.Stdout, "SYS: ", log.Ldate|log.Ltime|log.Llongfile)
	errLog = log.New(os.Stderr, "ERR: ", log.Ldate|log.Ltime|log.Llongfile)
}

func main() {
	initLogger()
	sysLog.Printf("Booting program...")
	// Root Endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, err := os.Hostname()
		if err != nil {
			errLog.Fatalf("Oops something went wrong while retrieving hostname")
		}
		sendEnvelope(w, 200, fmt.Sprintf("Hello folks! I am running on node: %s", host), nil)
	})
	// Healthcheck Endpoint
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		sendEnvelope(w, 200, fmt.Sprintf("PONG"), nil)
	})
	// Start a web server
	sysLog.Printf("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		errLog.Fatalf("Error starting server: %s", err)
	}
}
