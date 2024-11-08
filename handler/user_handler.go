package handler

import (
	"Notify/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key")) // Session store

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string `json:"name"`
		DomainID string `json:"domainid"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.Service.RegisterUser(req.Name, req.DomainID, req.Password, req.Role)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		DomainID string `json:"domainid"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.Service.LoginUser(req.DomainID, req.Password)
	if err != nil {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	// Create a session
	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Name
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60, // 1-month expiration
		HttpOnly: true,
	}

	if err := session.Save(r, w); err != nil {
		log.Println("Failed to save session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Login successful"))
}

func (h *UserHandler) AuthenticatedEndpoint(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.Values["user_id"].(string)
	username := session.Values["username"].(string)
	w.Write([]byte("Authenticated as " + username + " with ID: " + userID))
}
