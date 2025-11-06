# Minimal URL Shortener (Go + SQLite)

This is a lightweight URL shortener built using Go and SQLite.  
I created this project mainly for educational purposes — to understand how a URL shortener works internally, how routing and HTTP handlers function in Go, and how to persist data using SQLite.

---

## Project Goal

The goal of this project was to learn:

- How URL shorteners generate and store short codes.
- How HTTP request handling and JSON APIs work in Go.
- How to use SQLite for simple, persistent storage.
- How to build something end-to-end using only standard Go libraries.

This project was not built for production use, but as a way to better understand the backend process behind a common real-world application.

---

## Features

- Minimal HTML and JavaScript front-end.
- REST API for URL shortening.
- SQLite database for persistent data storage.
- Simple redirect functionality.
- Self-contained Go application (no external frameworks).

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|-------------|----------|
| Backend | Go (net/http) | Routing, API handling, serving HTML |
| Database | SQLite (`modernc.org/sqlite`) | Persistent storage |
| Frontend | HTML, CSS, JavaScript | Minimal UI |
| Architecture | Single-file monolith | Simple and easy to understand |

---

## Setup and Installation

1. Clone the repository
   ```bash
   git clone https://github.com/prityushsingh/GOurl.git
   cd GOurl
