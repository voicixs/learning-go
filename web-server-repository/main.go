package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// ============================================================================
// 1. DOMAIN MODEL
// ============================================================================

type User struct {
	Name string `json:"name"`
}

// ============================================================================
// 2. REPOSITORY INTERFACE (THE CONTRACT)
// ============================================================================
// Defines method signatures for data persistence without implementation details.
type UserRepository interface {
	Create(u User) int
	Get(id int) (User, bool)
	Delete(id int) bool
}

// ============================================================================
// 3. CONCRETE IMPLEMENTATION (IN-MEMORY MAP)
// ============================================================================
// Encapsulates the map and mutex inside a dedicated struct.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	store map[int]User
}

// Constructor function to initialize the repository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		store: make(map[int]User),
	}
}

// Create stores a new user and returns their generated ID.
// Uses a pointer receiver (*InMemoryUserRepository) to mutate state safely.
func (r *InMemoryUserRepository) Create(u User) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := len(r.store) + 1
	r.store[id] = u
	return id
}

// Get retrieves a user by ID.
func (r *InMemoryUserRepository) Get(id int) (User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.store[id]
	return u, exists
}

// Delete removes a user by ID and reports whether the user existed.
func (r *InMemoryUserRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return false
	}
	delete(r.store, id)
	return true
}

// ============================================================================
// 4. HTTP HANDLER STRUCT (DEPENDENCY INJECTION)
// ============================================================================
// UserHandler holds a dependency on the UserRepository INTERFACE,
// not the concrete InMemoryUserRepository struct.
type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// Delegates storage completely to the interface
	h.repo.Create(u)
	w.WriteHeader(http.StatusNoContent)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := h.repo.Get(id)
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deleted := h.repo.Delete(id)
	if !deleted {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// 5. SERVER INITIALIZATION & ROUTING
// ============================================================================

func main() {
	// 1. Initialize repository (data layer)
	userRepo := NewInMemoryUserRepository()

	// 2. Inject repository into handler (transport layer)
	userHandler := NewUserHandler(userRepo)

	// 3. Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
