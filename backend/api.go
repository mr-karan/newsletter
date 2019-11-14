package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
)

// Response represents the standardized API response struct.
type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

const (
	redisNamespace      = "newsletter"
	confirmKeyNamespace = redisNamespace + ":confirm"
)

// Subscription represents the parameters for unmarshalling
// create subscription endpoint
type Subscription struct {
	EmailID string `json:"email"`
}

// sendResponse is used to send success response based on format defined in Response
func sendResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	// Standard marshalled envelope for success.
	a := Response{
		Data:    data,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "Something unexpected happened", nil)
		return
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var app = r.Context().Value("app").(*App)
	index, err := app.fs.Read("/static/index.html")
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("error loading index from stuffed binary"), nil)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(index)))
	return
}

// Root endpoint
func handleAPIRoot(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, fmt.Sprintf("Welcome to newsletter subscription API"), nil)
	return
}

// Healthcheck endpoint
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, fmt.Sprintf("healthy"), nil)
	return
}

// New subscription endpoint
func handleNewSubscription(w http.ResponseWriter, r *http.Request) {
	var app = r.Context().Value("app").(*App)
	// decode request payload in a struct
	var sub Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		app.logger.Printf("Error while parsing body %s", err)
		sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unable to parse the request"), nil)
		return
	}
	if len(sub.EmailID) > 254 || !govalidator.IsEmail(sub.EmailID) {
		sendResponse(w, http.StatusBadRequest, fmt.Sprintf("Email ID: %s is not valid", sub.EmailID), nil)
		return
	}
	// TODO: Integrate Mailgun and send a confirmation email
	token, err := generateToken(32)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unable to generate token"), nil)
		return
	}
	// store token in cache
	conn := app.cachePool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", fmt.Sprintf("%s:%s", confirmKeyNamespace, token), fmt.Sprintf(sub.EmailID))
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unable to store token in cache"), nil)
		return
	}
	sendResponse(w, http.StatusOK, fmt.Sprintf(sub.EmailID), nil)
	return
}

// Confirm email endpoint

func handleConfirmEmail(w http.ResponseWriter, r *http.Request) {
	// var app = r.Context().Value("app").(*App)
}
