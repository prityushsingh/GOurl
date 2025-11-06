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

// 🧩 Home page (UI)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>GOurl</title>
		<style>
			body { font-family: Arial, sans-serif; text-align: center; padding: 50px; background-color: #f7f7f7; }
			h1 { color: #333; }
			input { padding: 10px; width: 300px; font-size: 16px; border-radius: 6px; border: 1px solid #ccc; }
			button { padding: 10px 20px; font-size: 16px; margin-left: 10px; border-radius: 6px; border: none; background-color: #007bff; color: white; cursor: pointer; }
			button:hover { background-color: #0056b3; }
			.result { margin-top: 20px; font-size: 18px; color: #007bff; }
			.error { color: red; margin-top: 20px; font-size: 16px; }
			.list { margin-top: 30px; font-size: 15px; text-align: left; display: inline-block; }
		</style>
	</head>
	<body>
		<h1>🔗 Minimal URL Shortener</h1>
		<input type="text" id="longUrl" placeholder="Enter your long URL here" />
		<button onclick="shorten()">Shorten</button>
		<div id="message"></div>
		<div class="list" id="list"></div>

		<script>
			async function shorten() {
				const url = document.getElementById('longUrl').value.trim();
				const message = document.getElementById('message');
				message.innerHTML = "";
				if (!url) {
					message.innerHTML = '<div class="error">Please enter a valid URL</div>';
					return;
				}

				try {
					const res = await fetch('/shorten', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ url })
					});
					if (!res.ok) {
						throw new Error('Failed to shorten URL');
					}
					const data = await res.json();
					message.innerHTML = '<div class="result">Shortened URL: <a href="' + data.short_url + '" target="_blank">' + data.short_url + '</a></div>';
					document.getElementById('longUrl').value = '';
					loadList();
				} catch (err) {
					message.innerHTML = '<div class="error">' + err.message + '</div>';
				}
			}

			async function loadList() {
				try {
					const res = await fetch('/list');
					const data = await res.json();
					let html = '<h3>🗂️ Recently Shortened URLs</h3><ul>';
					for (const [code, url] of Object.entries(data)) {
						html += '<li><b><a href="http://localhost:8080/' + code + '" target="_blank">' + code + '</a></b> → ' + url + '</li>';
					}
					html += '</ul>';
					document.getElementById('list').innerHTML = html;
				} catch {
					document.getElementById('list').innerHTML = '<p class="error">Could not load list</p>';
				}
			}

			loadList();
		</script>
	</body>
	</html>
	`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

// 🧩 POST /shorten — Create a new short URL
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

// 🧩 GET /<short_code> — Redirect
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // remove leading "/"
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

// 🧩 GET /list — View all URLs
func listHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT short_code, original_url FROM urls ORDER BY id DESC LIMIT 10")
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

// ======== Main ========

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/list", listHandler)

	// Single "/" dispatcher: home at "/", redirect for everything else
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			homeHandler(w, r)
		case "/favicon.ico", "/robots.txt":
			// no-op
		default:
			redirectHandler(w, r)
		}
	})

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
