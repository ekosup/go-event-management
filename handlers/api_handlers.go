package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"event-management/db"
	"event-management/models"
	"event-management/services"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest is the request body for user registration
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ReturnTo string `json:"return_to"`
}

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := models.User{Name: req.Name, Email: req.Email, Password: string(hashedPassword)}
	if result := db.DB.Create(&user); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateEvent: Handling event creation request.")
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("CreateEvent: Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uint)
	event.UserID = userID
	log.Printf("CreateEvent: UserID from context: %v\n", userID)

	if result := db.DB.Create(&event); result.Error != nil {
		log.Printf("CreateEvent: Error saving event to database: %v\n", result.Error)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("CreateEvent: Event created successfully.")
	w.WriteHeader(http.StatusCreated)
}

func CreateGuest(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateGuest: Handling guest creation request.")
	var guest models.Guest
	if err := json.NewDecoder(r.Body).Decode(&guest); err != nil {
		log.Printf("CreateGuest: Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("CreateGuest: Guest data received: %+v\n", guest)

	// Save the guest first to get an ID
	if result := db.DB.Create(&guest); result.Error != nil {
		log.Printf("CreateGuest: Error saving guest to database: %v\n", result.Error)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Generate QR code after guest has an ID
	qrCodeData := fmt.Sprintf("{\"guest_id\":%d, \"event_id\":%d}", guest.ID, guest.EventID)
	qrCode, err := services.GenerateQRCode(qrCodeData)
	if err != nil {
		log.Printf("CreateGuest: Failed to generate QR code: %v\n", err)
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}
	guest.QRCode = qrCode
	log.Printf("CreateGuest: QR Code generated: %s\n", qrCode)

	// Update the guest with the QR code
	if result := db.DB.Save(&guest); result.Error != nil {
		log.Printf("CreateGuest: Error updating guest with QR code: %v\n", result.Error)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("CreateGuest: Guest created and QR code updated successfully.")
	w.WriteHeader(http.StatusCreated)
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("GetEvents: Handling event retrieval request.")
	userID := r.Context().Value("userID").(uint)
	log.Printf("GetEvents: UserID from context: %v\n", userID)

	var events []models.Event
	if result := db.DB.Where("user_id = ?", userID).Find(&events); result.Error != nil {
		log.Printf("GetEvents: Error fetching events from database: %v\n", result.Error)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GetEvents: Found %d events for UserID %v.\n", len(events), userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("GetEvent: Handling event retrieval request.")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/events/")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("GetEvent: Invalid event ID: %v\n", err)
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uint)

	var event models.Event
	if result := db.DB.Preload("Guests").Where("id = ? AND user_id = ?", uint(id), userID).First(&event); result.Error != nil {
		log.Printf("GetEvent: Error fetching event from database: %v\n", result.Error)
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	log.Printf("GetEvent: Found event %v for UserID %v.\n", event.ID, userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateEvent: Handling event update request.")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/events/")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("UpdateEvent: Invalid event ID: %v\n", err)
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uint)

	var event models.Event
	if result := db.DB.Where("id = ? AND user_id = ?", uint(id), userID).First(&event); result.Error != nil {
		log.Printf("UpdateEvent: Error fetching event from database: %v\n", result.Error)
		http.Error(w, "Event not found or unauthorized", http.StatusNotFound)
		return
	}

	var updatedEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
		log.Printf("UpdateEvent: Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.DB.Model(&event).Updates(updatedEvent)
	log.Printf("UpdateEvent: Event %v updated successfully.\n", event.ID)
	w.WriteHeader(http.StatusOK)
}

// LogoutUser handles user logout
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

// LoginUser handles user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if result := db.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key")) // Use secret from config
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Path:     "/",
	})

	redirectURL := "/"
	if req.ReturnTo != "" {
		redirectURL = req.ReturnTo
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"redirect": redirectURL})
}
