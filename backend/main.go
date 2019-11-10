package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/knadh/stuffbin"
)

var (
	version = "1.0.1"
	sysLog  *log.Logger
	errLog  *log.Logger
)

var regexEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Response represents the standardized API response struct.
type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

// Subscription represents the parameters for unmarshalling
// create subscription endpoint
type Subscription struct {
	EmailID string `json:"email"`
}

// sendEnvelope is used to send success response based on format defined in Response
func sendEnvelope(w http.ResponseWriter, code int, message string, data interface{}) {
	// Standard marshalled envelope for success.
	a := Response{
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

func initFS() {
	fs, err := stuffbin.UnStuff(os.Args[0])
	if err != nil {
		errLog.Fatalf("error reading stuffed binary: %v", err)
	}
	fmt.Println("loaded files", fs.List())
}

func main() {
	initLogger()
	initFS()
	sysLog.Printf("Booting program...")
	// Root Endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sendEnvelope(w, http.StatusOK, fmt.Sprintf("Welcome to newsletter subscription API"), nil)
		return
	})
	// Healthcheck endpoint.
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		sendEnvelope(w, http.StatusOK, fmt.Sprintf("PONG"), nil)
		return
	})
	// Create subscription endpoint.
	http.HandleFunc("/api/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			sendEnvelope(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s request is not allowed", r.Method), nil)
			return
		}
		// decode request payload in a struct
		var sub Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			errLog.Printf("Error while parsing body %s", err)
			sendEnvelope(w, http.StatusInternalServerError, fmt.Sprintf("Unable to parse the request"), nil)
			return
		}
		if len(sub.EmailID) > 254 || !regexEmail.MatchString(sub.EmailID) {
			sendEnvelope(w, http.StatusBadRequest, fmt.Sprintf("EMail ID: %s is not valid", sub.EmailID), nil)
			return
		}
		sendEnvelope(w, http.StatusOK, fmt.Sprintf(sub.EmailID), nil)
		return
	})
	// Confirm email endpoint.
	http.HandleFunc("/api/confirmation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			sendEnvelope(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s request is not allowed", r.Method), nil)
			return
		}
		// decode request payload in a struct
		sendEnvelope(w, http.StatusOK, "wip", nil)
		return
	})
	// Start a web server
	sysLog.Printf("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		errLog.Fatalf("Error starting server: %s", err)
	}
}
