# NaRakhate - Betting Website

**NaRakhate** is a betting website project designed to provide users with the ability to manage bets, events, and categories. It includes server-side logic, database interactions, and an API for handling HTTP requests.

---

## Key Features
- User registration and authentication.
- Event and category management.
- Creation and management of bets.
- Logging for monitoring operations.
- Database integration for storing and processing data.
- Docker-based deployment and management.

---

## Project Structure

### 1. `config/`
Stores the initial application configuration, including:
- Database connection settings.
- API and server parameters.

### 2. `data/`
Folder for storing logs:
- Application logs are written to `data/app.log`.

### 3. `flags/`
Contains flags used for configuring the application.

### 4. `internal/`
The core logic of the project:
- **API:** Handles HTTP requests, including authentication, bets, events, and categories.
- **Business Logic:** Implements key betting website features.
- **Database (DB):** Handles data operations and connections.
- **Models:** Data structures for users, bets, events, and categories.
- **Server:** Manages server initialization and configurations.

### 5. `logging/`
Logging module:
- Tracks user actions and system events.
- Logs errors and information into `data/app.log`.

### 6. `utils/`
Utility functions to simplify routine tasks.

### 7. `docker-compose.yml`
Manages containerized application deployment with Docker Compose.

### 8. `Dockerfile`
Builds the Docker image for the application.

### 9. `go.mod` and `go.sum`
Files managing Go dependencies:
- `go.mod`: Defines module dependencies.
- `go.sum`: Stores checksums for module versions.

### 10. `init.sql`
SQL scripts for:
- Creating database tables.
- Importing initial data (e.g., categories, events).

---

## Installation

### 1. Clone the repository
Clone the project to your local machine:
```bash
git clone https://github.com/suyundykovv/NaRakhate.git
cd NaRakhate
```

### 2. Install dependencies
Ensure **Go** 1.18 or higher is installed, then run:
```bash
go mod tidy
```

### 3. Set up the database
Run the SQL scripts from `init.sql` to create tables and import initial data into your database.

### 4. Run using Docker
If you're using Docker, start the application with:
```bash
docker-compose up
```

### 5. Run locally
To run the server without Docker:
```bash
go run main.go
```

The server will be available at:
```text
http://localhost:8080
```

---

## API
[Add details about API endpoints, e.g.:]
- **POST /login** — User login.
- **GET /events** — Retrieve a list of events.
- **POST /bet** — Create a new bet.
- **GET /categories** — Retrieve a list of categories.

---

## Logging
All application activities are logged in `data/app.log`. The logging includes:
- Errors.
- Key user actions.
- Internal system operations.

---

## Contacts
- **Author**: Suyundykovv
- **GitHub**: [https://github.com/suyundykovv](https://github.com/suyundykovv)