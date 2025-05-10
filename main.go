package main

import (
	"database/sql"
	"encoding/json"
	_ "event-manager/docs" // This will be generated
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Event Manager API
// @version 1.0
// @description This is a sample event management server.
// @host localhost:8081
// @BasePath /
type Event struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Team Meeting"`
	Description string    `json:"description,omitempty" example:"Weekly team sync"`
	Location    string    `json:"location" example:"Conference Room A"`
	StartTime   time.Time `json:"start_time" example:"2024-03-20T10:00:00Z"`
	EndTime     time.Time `json:"end_time" example:"2024-03-20T11:00:00Z"`
	CreatedBy   string    `json:"created_by" example:"john.doe"`
	CreatedAt   time.Time `json:"created_at" example:"2024-03-19T15:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-03-19T15:00:00Z"`
}

var db *sql.DB

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int         `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResponse represents a validation error response
type ValidationResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	var err error
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "event_management")

	log.Printf("Connecting to database at %s:%s...", dbHost, dbPort)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Println("dsn: ", dsn)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Testing database connection...")
	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize database schema
	log.Println("Initializing database schema...")
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database schema initialized successfully")

	// Swagger documentation endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/events/", eventHandler)
	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func initDB() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			location VARCHAR(255) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			created_by VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listEvent(w, r)
	case http.MethodPost:
		createEvent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	fmt.Println("header: ", w.Header().Get("Content-Type"))
}

// @Summary List all events
// @Description Get paginated list of events from the database
// @Tags events
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Success 200 {object} PaginatedResponse{data=[]Event}
// @Router /events [get]
func listEvent(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query string
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 15
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count of events
	var totalItems int
	err = db.QueryRow("SELECT COUNT(*) FROM events").Scan(&totalItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalItems + pageSize - 1) / pageSize

	// Get paginated events
	rows, err := db.Query(`
		SELECT id, title, description, location, start_time, end_time, created_by, created_at, updated_at 
		FROM events 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	// Create paginated response
	response := PaginatedResponse{
		Data:       events,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Create a new event
// @Description Create a new event in the database
// @Tags events
// @Accept json
// @Produce json
// @Param event body Event true "Event object"
// @Success 201 {object} Response
// @Failure 422 {object} ValidationResponse
// @Failure 500 {string} string
// @Router /events [post]
func createEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the event
	errors := validateEvent(event)
	if len(errors) > 0 {
		validationResponse := ValidationResponse{
			Message: "Validation errors",
			Errors:  errors,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(validationResponse)
		return
	}

	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	var id int
	err := db.QueryRow(`
		INSERT INTO events (title, description, location, start_time, end_time, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		event.Title, event.Description, event.Location, event.StartTime, event.EndTime,
		event.CreatedBy, event.CreatedAt, event.UpdatedAt).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event.ID = id
	resp := Response{
		Message: "Event created successfully",
		Data:    event,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/events/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	switch r.Method {
	case http.MethodGet:
		getEvent(w, id)
	case http.MethodPut:
		updateEvent(w, id, r)
	case http.MethodDelete:
		deleteEvent(w, id)
	default:
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
}

// @Summary Get a specific event
// @Description Get an event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} Event
// @Router /events/{id} [get]
func getEvent(w http.ResponseWriter, id int) {
	e, err := getEventFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(e)
}

// @Summary Update an event
// @Description Update an existing event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body Event true "Event object"
// @Success 200 {object} Response
// @Router /events/{id} [put]
func updateEvent(w http.ResponseWriter, id int, r *http.Request) {
	event, err := getEventFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var updatedEvent Event
	if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updatedEvent.ID = id
	updatedEvent.CreatedAt = event.CreatedAt
	updatedEvent.UpdatedAt = time.Now()
	_, err = db.Exec(`
			UPDATE events SET title=$1, description=$2, location=$3, start_time=$4, end_time=$5, created_by=$6, updated_at=$7
			WHERE id=$8`,
		updatedEvent.Title, updatedEvent.Description, updatedEvent.Location, updatedEvent.StartTime, updatedEvent.EndTime, updatedEvent.CreatedBy, updatedEvent.UpdatedAt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := Response{
		Message: "event updated",
		Data:    updatedEvent,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} Response
// @Router /events/{id} [delete]
func deleteEvent(w http.ResponseWriter, id int) {
	result, err := db.Exec(`DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "No event found with the given id", http.StatusNotFound)
		return
	}

	resp := Response{
		Message: "Event deleted successfully",
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getEventFromDb(id int) (*Event, error) {
	var e Event
	err := db.QueryRow(`Select * from events where id = $1`, id).Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// validateEvent performs validation on the event data
func validateEvent(event Event) []ValidationError {
	var errors []ValidationError

	// Validate required fields
	if event.Title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title is required",
		})
	}

	if event.Location == "" {
		errors = append(errors, ValidationError{
			Field:   "location",
			Message: "location is required",
		})
	}

	if event.CreatedBy == "" {
		errors = append(errors, ValidationError{
			Field:   "created_by",
			Message: "created_by is required",
		})
	} else if !isValidEmail(event.CreatedBy) {
		errors = append(errors, ValidationError{
			Field:   "created_by",
			Message: "created_by must be a valid email address",
		})
	}

	// Validate time fields
	now := time.Now()
	if event.StartTime.IsZero() {
		errors = append(errors, ValidationError{
			Field:   "start_time",
			Message: "start_time is required",
		})
	} else if event.StartTime.Before(now) {
		errors = append(errors, ValidationError{
			Field:   "start_time",
			Message: "start_time must be in the future",
		})
	}

	if event.EndTime.IsZero() {
		errors = append(errors, ValidationError{
			Field:   "end_time",
			Message: "end_time is required",
		})
	} else if event.EndTime.Before(event.StartTime) {
		errors = append(errors, ValidationError{
			Field:   "end_time",
			Message: "end_time must be after start_time",
		})
	}

	return errors
}

// isValidEmail checks if a string is a valid email address
func isValidEmail(email string) bool {
	// Basic email validation
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := strings.Index(email, "@")
	if at == -1 || at == 0 || at == len(email)-1 {
		return false
	}
	dot := strings.LastIndex(email[at:], ".")
	if dot == -1 || dot == 0 || dot == len(email[at:])-1 {
		return false
	}
	return true
}
