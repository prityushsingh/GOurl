# 🔗 Minimal URL Shortener (Go + SQLite)

A lightweight, self-contained **URL Shortener** built with **Go** and **SQLite**.  
This project was created primarily for **educational purposes** — to understand how URL shorteners work under the hood,  
including concepts like routing, HTTP handlers, database persistence, and server-side HTML rendering.

---

## 🎯 Project Goal

The main goal of this project was to:
- Learn how backend systems like URL shorteners generate and store short codes.
- Understand HTTP request handling, routing, and JSON APIs in Go.
- Work with SQLite to add persistent data storage to a Go web app.
- Build something functional from scratch — not just follow a tutorial.

It’s a small project, but it helped me grasp how a simple idea (shortening URLs) combines multiple backend concepts together.

---

## 🚀 Features

- ✨ Minimal, clean front-end interface (pure HTML + JS)
- 🧠 REST API for programmatic URL shortening
- 🗃️ SQLite database for persistent storage
- 🔁 Instant redirect to original URLs
- 💡 All-in-one Go application (no external server or front-end build tools)

---

## 🧠 Tech Stack

| Layer | Technology | Purpose |
|-------|-------------|----------|
| **Backend** | Go (net/http) | Handles routes, API, and serving HTML |
| **Database** | SQLite (via `modernc.org/sqlite`) | Stores short codes and original URLs |
| **Frontend** | HTML, CSS, JavaScript | Provides the minimal UI |
| **Architecture** | Single-file monolith | Simple structure for learning |

---

## ⚙️ Installation & Setup

### 1️⃣ Clone the repository
```bash
git clone https://github.com/<your-username>/GOurl.git
cd GOurl
