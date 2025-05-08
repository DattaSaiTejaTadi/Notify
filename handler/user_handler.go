package handler

import (
	"Notify/constants"
	"Notify/middelware"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))

var userConnections = make(map[string]chan string) // key: userID, value: channel for SSE
var mu sync.Mutex

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	name := r.FormValue("name")
	domainid := r.FormValue("domainid")
	password := r.FormValue("password")
	role := r.FormValue("role")
	user, err := h.service.RegisterUser(name, domainid, password, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	//fmt.Print("got the request")
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	domainid := r.FormValue("domainid")
	password := r.FormValue("password")
	user, err := h.service.LoginUser(domainid, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	jwtDetails := map[string]interface{}{
		constants.UserID:   user.ID,
		constants.Role:     user.Role,
		constants.UserName: user.Name,
	}
	token, err := middelware.GenerateJWT(jwtDetails)
	if err != nil {
		http.Error(w, "Cannot pass token", http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:     constants.JWT_COOKIE,                // Name of the cookie
		Value:    token,                               // JWT token as the value
		HttpOnly: true,                                // Security flag to prevent JavaScript access
		Secure:   false,                               // Set to true for HTTPS
		SameSite: http.SameSiteStrictMode,             // CSRF protection: strict or lax
		Expires:  time.Now().Add(30 * 24 * time.Hour), // Expiration time (30 days)
	}
	http.SetCookie(w, cookie)
	if user.Role == "senior" {
		http.Redirect(w, r, "/seniorHome", http.StatusPermanentRedirect)
		//w.Write([]byte("Login successful")))
	}
	if user.Role == "member" {
		http.Redirect(w, r, "/memberHome", http.StatusPermanentRedirect)
		//w.Write([]byte("Login successful")))
	}
	http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
	//w.Write([]byte("Login successful"))
}

func (h *handler) AuthenticatedEndpoint(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	userDetails, ok := r.Context().Value(constants.UserDetailsFromToken).(map[string]interface{})
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	jsdata, _ := json.Marshal(userDetails)
	w.Write([]byte(jsdata))
}
func (h *handler) GetNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Http Method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Straming not supported", http.StatusInternalServerError)
	}

	userDetails, ok := r.Context().Value(constants.UserDetailsFromToken).(map[string]interface{})
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID := userDetails[constants.UserID].(string)
	err := h.service.UpdateActivityAsTrue(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Create a new channel to send events to this user
	messageChannel := make(chan string)
	// Store the message channel for the user
	mu.Lock()
	userConnections[userID] = messageChannel
	mu.Unlock()
	// Use a channel to monitor for client disconnect
	quit := make(chan bool)

	// Start a goroutine to detect client disconnect
	go func() {
		select {
		case <-r.Context().Done():
			// The client disconnected
			fmt.Println("Client disconnected")
			quit <- true
		}
	}()
	for {
		select {
		case msg, ok := <-messageChannel: // If we receive a message for this user
			if !ok { // Check if the channel was closed
				log.Printf("Channel for user %s closed", userID)
				return
			}
			// Send the received message as an SSE event to the user
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush() // Flush the message to the client

		case <-quit:
			_ = h.service.UpdateActivityAsFalse(userID)

		}
	}
}

// Function to send a message to a specific user
func sendMessageToUser(userID, message string) {
	mu.Lock()
	defer mu.Unlock()

	// Check if the user has an active connection
	if ch, exists := userConnections[userID]; exists {
		ch <- message // Send the message to the user's channel
	} else {
		log.Printf("No active connection for user %s\n", userID)
	}
}

// New handler to trigger event manually by sending a message to a user
func (h *handler) Notify(w http.ResponseWriter, r *http.Request) {
	// Only allow Post method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var request struct {
		UserID  string `json:"userID"`
		Message string `json:"message"`
	}

	// Decode the JSON request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Send the message to the specified user
	sendMessageToUser(request.UserID, request.Message)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message sent to user %s", request.UserID)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
	// Clear the JWT cookie by setting its value to an empty string
	http.SetCookie(w, &http.Cookie{
		Name:    constants.JWT_COOKIE,
		Value:   "",              // Empty value to clear it
		Path:    "/",             // Path where the cookie is valid, typically "/"
		Expires: time.Unix(0, 0), // Set expiration to a time in the past
	})
	// Optionally, send a message to the client that the session has been deleted
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (h *handler) GetMembers(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// Check if flusher is available
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	//Fetch and send initial data
	users, err := h.service.GetMembers()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	userResponse, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", userResponse)
	flusher.Flush() // Immediately flush the initial response

}

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(users)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	err := h.service.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

// Status check
func (h *handler) StatusCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
	status := r.URL.Query().Get("status")
	log.Println("status", status)
	w.WriteHeader(http.StatusOK)
}
