# Event Management API (Go + MySQL)

A simple RESTful API built with Go's `net/http` and `database/sql` packages for managing events.

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Database**: MySQL
- **Packages Used**:
    - `github.com/go-sql-driver/mysql`
    - `encoding/json`
    - `net/http`
    - `database/sql`

---

## 🗄 Database Setup

1. Create a MySQL database called `event_managment` (or modify the DSN in `main.go`).
2. Run the following SQL to create the `events` table:

```sql
CREATE TABLE events (
  id INT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  location VARCHAR(255),
  start_time DATETIME,
  end_time DATETIME,
  created_by VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME
);
```

---

## 🚀 Getting Started

1. Clone the repo and open the project directory.
2. Update your MySQL connection in `main.go`:
   ```go
   dsn := "root:root@tcp(127.0.0.1:3306)/event_management?parseTime=true"
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
4. Server will run at:  
   `http://localhost:8081`

---

## 📦 API Endpoints

### `GET /events`
Returns all events.

### `POST /events`
Create a new event.  
**Request Body (JSON):**
```json
{
  "title": "Meeting",
  "description": "Project kickoff",
  "location": "Office",
  "start_time": "2025-05-08T10:00:00Z",
  "end_time": "2025-05-08T11:00:00Z",
  "created_by": "alice@example.com"
}
```

### `GET /events/{id}`
Get a single event by ID.

### `PUT /events/{id}`
Update an event.  
**Request Body:** Same as POST.

### `DELETE /events/{id}`
Delete an event by ID.

---

## ⚠ Notes

- `Content-Type` must be `application/json` for POST and PUT.
- Time format must follow RFC3339: `"2025-05-08T10:00:00Z"`.
- Errors are returned as plain text or JSON.

---