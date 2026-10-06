package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// User struct defines the user entity with JSON tags for serialization.
type User struct {
	Name string `json:"name"`
}

// Global in-memory cache acting as a simulated database, along with an RWMutex for thread safety.
var (
	userCache  = make(map[int]User)
	cacheMutex sync.RWMutex
)

// handleRoot handles requests to the root "/" endpoint.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

// createUser parses JSON input and inserts a new user into the in-memory database.
func createUser(w http.ResponseWriter, r *http.Request) {
	var u User

	// Decode incoming JSON body into User struct
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if u.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// Lock mutex for writing to map safely across goroutines
	cacheMutex.Lock()
	id := len(userCache) + 1
	userCache[id] = u
	cacheMutex.Unlock()

	// Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// getUser retrieves a user by ID using Go 1.22 PathValue routing.
func getUser(w http.ResponseWriter, r *http.Request) {
	// Extract path parameter '{id}' from request URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read-lock mutex for thread-safe map access
	cacheMutex.RLock()
	u, ok := userCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Convert User struct to JSON
	jsonResp, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content header and write JSON response with 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// deleteUser removes a user by ID from the in-memory database.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheMutex.Lock()
	if _, ok := userCache[id]; !ok {
		cacheMutex.Unlock()
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	delete(userCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Example CRUD calls (run the server first: `go run .`, then copy/paste)
//
// Root / health check:
//
//	curl -i http://localhost:8080/
//
// CREATE a user (expects JSON body with "name"; returns 204 No Content):
//
//	curl -i -X POST http://localhost:8080/users \
//	  -H "Content-Type: application/json" \
//	  -d '{"name":"Alice"}'
//
//	Create a second user so there are a couple of IDs to play with:
//	curl -i -X POST http://localhost:8080/users \
//	  -H "Content-Type: application/json" \
//	  -d '{"name":"Bob"}'
//
//	Validation error (missing name -> 400 Bad Request):
//	curl -i -X POST http://localhost:8080/users \
//	  -H "Content-Type: application/json" \
//	  -d '{}'
//
// READ a user by ID (IDs start at 1; returns JSON with 200 OK):
//
//	curl -i http://localhost:8080/users/1
//
//	Not found (-> 404):
//	curl -i http://localhost:8080/users/999
//
//	Bad ID (-> 400):
//	curl -i http://localhost:8080/users/abc
//
// UPDATE: not implemented yet (no PUT/PATCH route). A future route could be:
//
//	curl -i -X PUT http://localhost:8080/users/1 \
//	  -H "Content-Type: application/json" \
//	  -d '{"name":"Alice Updated"}'
//
// DELETE a user by ID (returns 204 No Content; 404 if missing):
//
//	curl -i -X DELETE http://localhost:8080/users/1
//
//	Verify it is gone (-> 404):
//	curl -i http://localhost:8080/users/1
//
// Note: data lives in memory only and is reset every time the server restarts.
// ---------------------------------------------------------------------------
func main() {
	// Create a new request multiplexer (ServeMux)
	mux := http.NewServeMux()

	// Register routes with HTTP method matching and path variables
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	fmt.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}
