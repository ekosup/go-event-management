package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"event-management/config"
	"event-management/db"
	"event-management/handlers"
	"event-management/middleware"
)

func main() {
	cfg := config.LoadConfig()

	db.ConnectDatabase(cfg)

	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	homepageHandler := http.HandlerFunc(handlers.ShowHomepage)
	mux.Handle("/", middleware.AuthMiddleware(homepageHandler))

	// Page handlers
	mux.HandleFunc("/login", handlers.ShowLoginPage)
	mux.HandleFunc("/register", handlers.ShowRegisterPage)

	mux.HandleFunc("/logout", handlers.LogoutUser)

	mux.HandleFunc("/api-test", handlers.ShowApiTestPage)

	// API handlers
	mux.HandleFunc("/api/register", handlers.RegisterUser)
	mux.HandleFunc("/api/login", handlers.LoginUser)

	eventDetailsHandler := http.HandlerFunc(handlers.ShowEventDetailsPage)
	mux.Handle("/events/", middleware.AuthMiddleware(eventDetailsHandler))
	dashboardHandler := http.HandlerFunc(handlers.ShowDashboardPage)
	mux.Handle("/dashboard", middleware.AuthMiddleware(dashboardHandler))

	eventsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.CreateEvent(w, r)
		} else if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/events/") && len(strings.TrimPrefix(r.URL.Path, "/api/events/")) > 0 {
			handlers.GetEvent(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetEvents(w, r)
		} else if r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/events/") && len(strings.TrimPrefix(r.URL.Path, "/api/events/")) > 0 {
			handlers.UpdateEvent(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.Handle("/api/events/", middleware.AuthMiddleware(eventsHandler))

	createGuestHandler := http.HandlerFunc(handlers.CreateGuest)
	mux.Handle("/api/guests", middleware.AuthMiddleware(createGuestHandler))

	log.Printf("Server starting on port %s\n", cfg.ServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
