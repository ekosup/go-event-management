# **Role**

Act as an expert full-stack developer specializing in Go (Golang). Your task is to generate a complete, production-ready web application based on the provided system design document, using Go for both the backend and frontend HTML templating.

# **Task**

Create a monolithic web application for a "QR Code Presence & Event Management" system. The application must be written in Go (version 1.23) and should serve both a JSON API and a server-rendered HTML frontend. It should strictly adhere to the features and data models defined in the system design.

# **Core Requirements & Technologies**

### **Backend**

* **Language:** Go 1.23  
* **Web Server/Router:** net/http (standard library)  
* **Database:** PostgreSQL  
* **ORM:** GORM (gorm.io/gorm and gorm.io/driver/postgres)  
* **Authentication:** JWT (github.com/golang-jwt/jwt/v5)  
* **Password Hashing:** golang.org/x/crypto/bcrypt  
* **QR Code Generation:** github.com/skip2/go-qrcode  
* **Configuration:** Use environment variables (e.g., using github.com/joho/godotenv for local development).

### **Frontend**

* **Templating:** html/template (Go standard library)  
* **Styling:** Tailwind CSS  
* **Interactivity:** Alpine.js  
* **Dependencies:** Manage frontend dependencies with package.json (npm/yarn/pnpm).

# **Project Structure**

Organize the code into a clean, scalable monolithic structure. Please generate separate code blocks for each file.  
/presence-app  
|-- go.mod  
|-- main.go                 \# Application entry point, router setup  
|-- .env.example            \# Example environment variables  
|-- package.json            \# Frontend dependencies (Tailwind, Alpine.js)  
|-- tailwind.config.js      \# Tailwind configuration  
|  
|-- /config/                \# Configuration loading  
|-- /db/                    \# Database connection  
|-- /models/                \# GORM models  
|-- /handlers/              \# HTTP handlers for both pages and API endpoints  
|-- /middleware/            \# net/http compatible middleware  
|-- /services/              \# Business logic (QR, Email)  
|  
|-- /templates/             \# Go HTML templates  
|   |-- /layouts/  
|   |   |-- base.html  
|   |-- /pages/  
|   |   |-- login.html  
|   |   |-- register.html  
|   |   |-- dashboard.html  
|   |-- /partials/  
|   |   |-- header.html  
|  
|-- /static/                \# Compiled frontend assets  
|   |-- css/  
|   |   |-- style.css       \# Output of Tailwind CSS  
|   |-- js/  
|  
|-- /assets/                \# Source frontend files  
    |-- css/  
    |   |-- input.css       \# Main CSS file for Tailwind directives

# **Detailed Implementation Steps**

### **1\. Backend Implementation (Models, DB, Config)**

* Translate the SQL schema into GORM model structs (user.go, event.go, guest.go).  
* Implement ConnectDatabase() using GORM and credentials from environment variables.  
* Implement configuration loading for DB\_\*, JWT\_SECRET\_KEY, SERVER\_PORT.

### **2\. Frontend Setup**

* **package.json:** Define scripts to build the CSS. Include tailwindcss and alpinejs as dev dependencies.  
  * A build:css script should run: tailwindcss \-i ./assets/css/input.css \-o ./static/css/style.css \--watch  
* **tailwind.config.js:** Configure Tailwind to scan .html files in the /templates directory for classes.  
* **assets/css/input.css:** Include the base Tailwind directives (@tailwind base; @tailwind components; @tailwind utilities;).

### **3\. Web Server & Routing (main.go)**

* Use the standard net/http package. Create a new http.ServeMux router.  
* **Static Files:** Register a file server to serve the /static/ directory.  
* **Routing:**  
  * Define routes for serving HTML pages (e.g., /login, /dashboard). These handlers will parse and execute Go templates.  
  * Define API endpoint routes (e.g., POST /api/login). These will handle form submissions from the frontend (via Alpine.js), process data, and return JSON.  
* Start the server with http.ListenAndServe.

### **4\. Middleware (/middleware)**

* Create net/http compatible middleware. This will be a function that takes an http.Handler and returns an http.Handler.  
* **AuthMiddleware:**  
  1. Wrap API routes that require authentication.  
  2. Extract the JWT from a cookie or Authorization header.  
  3. Validate the token. If valid, add user info to the request context using context.WithValue.  
  4. If invalid, redirect to the login page or return a 401 Unauthorized for API requests.

### **5\. Handlers (/handlers)**

* **Page Handlers:** These functions will have the signature func(w http.ResponseWriter, r \*http.Request). They are responsible for:  
  * Fetching necessary data from the database.  
  * Parsing and executing the appropriate Go HTML templates from the /templates directory, passing the data to them.  
* **API Handlers:** These handlers will also have the http.HandlerFunc signature. They are responsible for:  
  * Decoding JSON request bodies.  
  * Performing business logic (e.g., creating a user, logging in, creating an event).  
  * Returning JSON responses to be consumed by Alpine.js on the frontend.  
  * The Login handler should set a secure, HttpOnly cookie containing the JWT.

### **6\. Templates & Frontend Interaction (/templates)**

* Create a base layout in /layouts/base.html that includes the \<head\> section (with the link to /static/css/style.css and the Alpine.js script tag) and defines blocks for content.  
* Page templates (/pages/\*.html) will extend this base layout.  
* Use Alpine.js within the HTML templates for client-side interactivity. For example, a login form would use x-data to hold form state and @submit.prevent to make a fetch call to the POST /api/login endpoint.

# **Output Format**

Please provide the complete, runnable Go code and frontend configuration files, each in a separate, clearly labeled markdown code block. Ensure the Go code is well-commented. Start with the package.json and tailwind.config.js files to establish the frontend environment.