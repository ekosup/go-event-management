package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Render a template
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(filepath.Join("templates", "layouts", "base.html"), filepath.Join("templates", "pages", tmpl), filepath.Join("templates", "partials", "header.html"))
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println(err)
	}
}

func ShowHomepage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "homepage.html", nil)
}

// ShowLoginPage renders the login page
func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
}

// ShowRegisterPage renders the register page
func ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register.html", nil)
}

func ShowApiTestPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "api_test.html", nil)
}

// ShowEventDetailsPage renders the event details page
func ShowEventDetailsPage(w http.ResponseWriter, r *http.Request) {
	log.Println("ShowEventDetailsPage: Rendering event details page.")
	renderTemplate(w, "event_details.html", nil)
}

// ShowDashboardPage renders the dashboard page
func ShowDashboardPage(w http.ResponseWriter, r *http.Request) {
	log.Println("ShowDashboardPage: Rendering dashboard page.")
	renderTemplate(w, "dashboard.html", nil)
}
