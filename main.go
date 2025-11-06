package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

// ======== Structs for API Request/Response ========

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// ======== Database Setup ========

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "data.db")
	if err != nil {
		panic(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_code TEXT NOT NULL UNIQUE,
		original_url TEXT NOT NULL
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Database initialized successfully")
}

// ======== Utility Function ========

func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// ======== Handlers ========

// POST /shorten
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode(6)

	// Store in DB
	_, err := db.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", shortCode, req.URL)
	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	response := ShortenResponse{
		ShortURL: fmt.Sprintf("http://localhost:8080/%s", shortCode),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("📦 Stored: %s → %s\n", shortCode, req.URL)
}

// GET /<short_code>
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // remove the leading "/"
	if code == "" {
		http.Error(w, "Missing short code", http.StatusBadRequest)
		return
	}

	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", code).Scan(&originalURL)
	if err == sql.ErrNoRows {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
	fmt.Printf("🔁 Redirected: %s → %s\n", code, originalURL)
}

// GET /list → show all stored URLs (for testing)
func listHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT short_code, original_url FROM urls")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var code, url string
		rows.Scan(&code, &url)
		result[code] = url
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ======== Main Function ========

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
