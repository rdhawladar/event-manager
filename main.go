package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var db *sql.DB

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/event_managment?parseTime=true" //change this by your mysql user and pass
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect")
	}
	if err := db.Ping(); err != nil {
		log.Fatal("db is unreachable")
	}
	/*db.SetConnMaxLifetime(1 * time.Minute)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(10)
	*/
	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/events/", eventHandler)
	fmt.Println("Server running at 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
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

func listEvent(w http.ResponseWriter, r *http.Request) {
	var events []Event
	rows, err := db.Query("Select * from events")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("header: ", w.Header().Get("Content-Type"))
	json.NewEncoder(w).Encode(events)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	result, err := db.Exec(`INSERT into events(title, description, location, start_time, end_time, created_by, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.CreatedBy, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	event.ID = int(id)
	resp := Response{
		Message: "Event created",
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
func getEvent(w http.ResponseWriter, id int) {
	e, err := getEventFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(e)
}

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
			UPDATE events SET title=?, description=?, location=?, start_time=?, end_time=?, created_by=?, updated_at=?
			WHERE id=?`,
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

func deleteEvent(w http.ResponseWriter, id int) {
	result, err := db.Exec(`Delete from events where id = ?`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		http.Error(w, "No event found with the given id", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Event deleted")
}

func getEventFromDb(id int) (*Event, error) {
	var e Event
	err := db.QueryRow(`Select * from events where id = ?`, id).Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
