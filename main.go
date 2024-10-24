package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-todo-app/handlers"
	"go-todo-app/models"
)

// LogEntry represents the structure of a log entry
type LogEntry struct {
	Timestamp    string `json:"timestamp"`
	Method       string `json:"method"`
	RequestURI   string `json:"request_uri"`
	StatusCode   int    `json:"status_code"`
	Duration     string `json:"duration"`
	ClientIP     string `json:"client_ip"`
	UserAgent    string `json:"user_agent"`
	ResponseSize int    `json:"response_size"`
}

// getIP extracts the client IP address from the RemoteAddr
func getIP(r *http.Request) string {
	// Check if the X-Forwarded-For header is set (contains a comma-separated list of IPs)
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		// Split by comma and take the first IP (the client IP)
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check if the X-Real-IP header is set
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr if headers are not set
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Fallback to full address if parsing fails
	}
	return ip
}

// LoggingMiddleware logs each HTTP request in JSON format with detailed information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code and response size
		crw := &CustomResponseWriter{w, http.StatusOK, 0}

		// Call the next handler
		next.ServeHTTP(crw, r)

		// Build log entry
		entry := LogEntry{
			Timestamp:    start.Format(time.RFC3339),
			Method:       r.Method,
			RequestURI:   r.RequestURI,
			StatusCode:   crw.StatusCode,
			Duration:     time.Since(start).String(),
			ClientIP:     getIP(r),
			UserAgent:    r.UserAgent(),
			ResponseSize: crw.Size,
		}

		// Log entry in JSON format
		logJSON, _ := json.Marshal(entry)
		fmt.Println(string(logJSON))
	})
}

// CustomResponseWriter is a wrapper around http.ResponseWriter to capture status code and response size
type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.StatusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}

func (crw *CustomResponseWriter) Write(data []byte) (int, error) {
	size, err := crw.ResponseWriter.Write(data)
	crw.Size += size
	return size, err
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println(`{"error":"DATABASE_URL not set"}`)
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Printf(`{"error":"Failed to connect to database: %v"}`, err)
		os.Exit(1)
	}

	db.AutoMigrate(&models.Todo{})

	r := mux.NewRouter()

	// Register API routes
	handlers.RegisterTodoHandlers(r, db)

	// Serve static files
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	// Wrap the router with the logging middleware
	loggedRouter := LoggingMiddleware(r)

	fmt.Println(`{"message":"Server started on :8080"}`)
	http.ListenAndServe("0.0.0.0:8080", loggedRouter)
}
