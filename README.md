# QR Code Presence & Event Management

This project is a monolithic web application for a "QR Code Presence & Event Management" system. The application is written in Go and serves both a JSON API and a server-rendered HTML frontend using Go's `html/template` package.

## Features

*   **User Management:** Secure user registration and login with JWT-based authentication.
*   **Event Creation:** Authenticated users can create and manage events.
*   **Guest Management:** Add guests to events and generate unique QR codes for each guest.
*   **QR Code Check-in:** A dedicated endpoint to validate a guest's QR code and mark them as present.
*   **Dashboard:** A user-facing dashboard to view and manage events and guests.
*   **RESTful API:** A JSON API for programmatic access to the application's features.

## Tech Stack

### Backend

*   **Language:** Go 1.23
*   **Web Server/Router:** `net/http` (standard library)
*   **Database:** PostgreSQL
*   **ORM:** GORM
*   **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`)
*   **Password Hashing:** `golang.org/x/crypto/bcrypt`
*   **QR Code Generation:** `github.com/skip2/go-qrcode`
*   **Configuration:** Environment variables (`github.com/joho/godotenv`)

### Frontend

*   **Templating:** `html/template` (Go standard library)
*   **Styling:** Tailwind CSS
*   **Interactivity:** Alpine.js

## Project Structure

```
/event-management
|-- go.mod
|-- main.go                 # Application entry point, router setup
|-- .env.example            # Example environment variables
|-- package.json            # Frontend dependencies
|-- tailwind.config.js      # Tailwind configuration
|
|-- /config/                # Configuration loading
|-- /db/                    # Database connection
|-- /models/                # GORM models
|-- /handlers/              # HTTP handlers for pages and API
|-- /middleware/            # net/http middleware
|-- /services/              # Business logic (QR, Email)
|
|-- /templates/             # Go HTML templates
|   |-- /layouts/
|   |-- /pages/
|   |-- /partials/
|
|-- /static/                # Compiled frontend assets
|   |-- css/
|
|-- /assets/                # Source frontend files
    |-- css/
```

## Getting Started

### Prerequisites

*   Go 1.23 or later
*   PostgreSQL
*   Node.js and npm (or yarn/pnpm)

### 1. Clone the Repository

```bash
git clone <repository-url>
cd presence-app
```

### 2. Configure Environment Variables

Create a `.env` file in the root of the project by copying the example file:

```bash
cp .env.example .env
```

Update the `.env` file with your database credentials and a secure JWT secret key:

```
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=5432
DB_SSLMODE=disable

JWT_SECRET_KEY=your-super-secret-key
SERVER_PORT=8080
```

### 3. Install Dependencies

**Backend (Go):**

```bash
go mod tidy
```

**Frontend (Node.js):**

```bash
npm install
```

### 4. Build Frontend Assets

This command will compile your Tailwind CSS. You can keep this running in a separate terminal during development to automatically rebuild on change.

```bash
npm run build:css
```

### 5. Run the Application

```bash
go run main.go
```

The application will be available at `http://localhost:8080`.

## API Endpoints

### Auth
*   `POST /api/register` - Register a new user.
*   `POST /api/login` - Login a user and receive a JWT cookie.

### Events (Protected)
*   `POST /api/events` - Create a new event.
*   `GET /api/events` - Get a list of events for the authenticated user.
*   `GET /api/events/{id}` - Get details for a specific event.

### Guests (Protected)
*   `POST /api/events/{id}/guests` - Add a guest to an event.
*   `GET /api/guests/scan/{qr_code}` - Scan a guest's QR code to check them in.

## License

This project is open-source and available under the [MIT License](LICENSE).
