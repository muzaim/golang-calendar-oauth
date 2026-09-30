# Golang Calendar & Authentication API

RESTful Web Service backend built with Golang, Gin Framework, GORM, MySQL, JWT Authentication, Google OAuth 2.0, and Google Calendar API integration.

---

## Tech Stack

* **Language**: Go (1.26+)
* **Framework**: Gin Gonic
* **Database**: MySQL 8.0
* **ORM**: GORM
* **Migration**: golang-migrate
* **Authentication**: JWT (Access Token + Refresh Token), Bcrypt
* **OAuth & Integration**: Google OAuth 2.0 & Google Calendar API v3
* **DevOps**: Docker, Docker Compose, Air (Hot Reload), Makefile

---

## Project Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── config/                  # Configuration loader (.env)
├── internal/
│   ├── dto/                 # Data transfer objects
│   ├── handler/             # HTTP request handlers
│   ├── middleware/          # JWT authentication middleware
│   ├── model/               # GORM database models
│   ├── repository/          # Data access layer
│   └── service/             # Business logic layer
├── routes/                  # Route definitions
├── database/
│   ├── database.go          # GORM MySQL connection
│   └── migrations/          # Migration SQL files
├── pkg/
│   ├── auth/                # JWT & Password hashing utilities
│   ├── oauth/               # Google OAuth & Calendar API client
│   └── response/            # Standardized API response helper
├── .env.example             # Environment variable template
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # Docker services configuration
├── Makefile                 # Development task runner commands
├── postman.json             # Postman collection for testing
└── go.mod                   # Go module dependencies
```

---

## Getting Started

### Prerequisites

* Go 1.22+
* Docker & Docker Compose
* MySQL 8.0 (if running locally without Docker)
* `golang-migrate` CLI (if running migrations manually)

---

### Step-by-Step Installation

#### 1. Clone Repository

```bash
git clone <repository_url>
cd golang-test
```

#### 2. Configure Environment Variables

Copy the `.env.example` template to `.env`:

```bash
cp .env.example .env
```

Edit `.env` and fill in your configuration:

```env
APP_PORT=8080
APP_ENV=development

DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=calendar_service

JWT_SECRET=supersecretjwtkeychangeinproduction
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_DAYS=7

GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

---

### Running the Application

#### Option A: Using Docker Compose (Recommended)

Start the MySQL database and application containers in one command:

```bash
docker compose up -d
```

View application logs:

```bash
docker compose logs -f app
```

Stop containers:

```bash
docker compose down
```

#### Option B: Local Development

1. Ensure MySQL is running and database `calendar_service` exists.
2. Run database migrations:
   ```bash
   make migrate-up
   ```
3. Start the application with Air hot-reload:
   ```bash
   make dev
   ```
   Or run standard build:
   ```bash
   make run
   ```

---

## Database Migrations

Use the Makefile commands to manage database migrations:

```bash
# Run all pending migrations
make migrate-up

# Rollback latest migration step
make migrate-down

# Create a new migration file
make migrate-create name=create_example_table

# Force set migration version
make migrate-force version=1
```

---

## Google OAuth 2.0 Setup

To enable Google Login and Google Calendar sync:

1. Open [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project and go to **APIs & Services > Credentials**.
3. Create **OAuth 2.0 Client IDs** (Application type: Web application).
4. Set Authorized redirect URIs to:
   `http://localhost:8080/api/auth/google/callback`
5. Enable **Google Calendar API** under **APIs & Services > Library**.
6. Copy `Client ID` and `Client Secret` into your `.env` file.

---

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Register new user | No |
| `POST` | `/api/auth/login` | Login user & return JWT tokens | No |
| `POST` | `/api/auth/refresh` | Refresh JWT access token | No |
| `POST` | `/api/auth/logout` | Revoke refresh token | No |
| `GET` | `/api/auth/me` | Get current user profile | Yes (Bearer Token) |
| `GET` | `/api/auth/google` | Redirect to Google OAuth Consent Screen | No |
| `GET` | `/api/auth/google/callback` | Google OAuth callback handler | No |

### Calendar Events

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/calendar/events` | List all user calendar events | Yes (Bearer Token) |
| `GET` | `/api/calendar/events/:id` | Get detail of a specific event | Yes (Bearer Token) |
| `POST` | `/api/calendar/events` | Create new event (syncs to Google) | Yes (Bearer Token) |
| `PUT` | `/api/calendar/events/:id` | Update event (syncs to Google) | Yes (Bearer Token) |
| `DELETE` | `/api/calendar/events/:id` | Delete event (syncs to Google) | Yes (Bearer Token) |

---

## Standard API Response Format

### Success Response

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "johndoe@example.com"
    }
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Invalid email or password",
  "errors": null
}
```

---

## Testing with Postman

1. Import `postman.json` into Postman.
2. Execute **Auth > Login** or **Auth > Google OAuth Callback**.
3. `access_token` and `refresh_token` will be set automatically as collection variables.
4. Execute any **Calendar Events** requests.

---

## Running Unit Tests

```bash
make test
```
