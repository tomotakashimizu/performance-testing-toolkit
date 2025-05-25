package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Item struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var items []Item
var nextID = 1

func main() {
	// 初期データを作成
	initializeData()

	r := mux.NewRouter()

	// ルートエンドポイント
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	// アイテム関連のエンドポイント
	r.HandleFunc("/api/v1/items", getItemsHandler).Methods("GET")
	r.HandleFunc("/api/v1/items", createItemHandler).Methods("POST")
	r.HandleFunc("/api/v1/items/{id}", getItemHandler).Methods("GET")
	r.HandleFunc("/api/v1/items/{id}", updateItemHandler).Methods("PUT")
	r.HandleFunc("/api/v1/items/{id}", deleteItemHandler).Methods("DELETE")

	// 負荷試験用の特別なエンドポイント
	r.HandleFunc("/api/v1/slow", slowHandler).Methods("GET")
	r.HandleFunc("/api/v1/random-delay", randomDelayHandler).Methods("GET")
	r.HandleFunc("/api/v1/cpu-intensive", cpuIntensiveHandler).Methods("GET")

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initializeData() {
	items = []Item{
		{ID: 1, Name: "Sample Item 1", Description: "This is a sample item", CreatedAt: time.Now()},
		{ID: 2, Name: "Sample Item 2", Description: "Another sample item", CreatedAt: time.Now()},
		{ID: 3, Name: "Sample Item 3", Description: "Yet another sample item", CreatedAt: time.Now()},
	}
	nextID = 4
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Welcome to Sample API",
		"version": "1.0.0",
		"status":  "running",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	item := Item{
		ID:          nextID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}
	items = append(items, item)
	nextID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, item := range items {
		if item.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for i, item := range items {
		if item.ID == id {
			items[i].Name = req.Name
			items[i].Description = req.Description
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(items[i])
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

// 負荷試験用のエンドポイント

func slowHandler(w http.ResponseWriter, r *http.Request) {
	// 固定で2秒の遅延
	time.Sleep(2 * time.Second)
	response := map[string]string{
		"message": "This endpoint is intentionally slow",
		"delay":   "2 seconds",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func randomDelayHandler(w http.ResponseWriter, r *http.Request) {
	// 0-1秒のランダムな遅延
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(delay)
	response := map[string]interface{}{
		"message": "This endpoint has random delay",
		"delay":   fmt.Sprintf("%.3f seconds", delay.Seconds()),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func cpuIntensiveHandler(w http.ResponseWriter, r *http.Request) {
	// CPU集約的な処理をシミュレート
	start := time.Now()
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}
	duration := time.Since(start)

	response := map[string]interface{}{
		"message":     "CPU intensive operation completed",
		"result":      sum,
		"duration_ms": duration.Milliseconds(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
