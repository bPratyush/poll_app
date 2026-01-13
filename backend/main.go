package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"poll_app/ent"
	"poll_app/handlers"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Get configuration from environment
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "postgres://localhost/poll_app?sslmode=disable")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	// Set JWT secret for handlers
	handlers.SetJWTSecret(jwtSecret)

	// Initialize database connection
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	// Create an ent driver from the sql.DB
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Initialize handlers
	h := handlers.NewHandler(client)

	// Setup router
	router := httprouter.New()

	// Auth routes
	router.POST("/api/auth/signup", h.SignUp)
	router.POST("/api/auth/login", h.Login)
	router.GET("/api/auth/me", h.AuthMiddleware(h.GetCurrentUser))

	// Poll routes
	router.GET("/api/polls", h.AuthMiddleware(h.ListPolls))
	router.POST("/api/polls", h.AuthMiddleware(h.CreatePoll))
	router.GET("/api/polls/:id", h.AuthMiddleware(h.GetPoll))
	router.PUT("/api/polls/:id", h.AuthMiddleware(h.UpdatePoll))
	router.DELETE("/api/polls/:id", h.AuthMiddleware(h.DeletePoll))

	// Vote routes
	router.POST("/api/polls/:id/vote", h.AuthMiddleware(h.Vote))
	router.DELETE("/api/polls/:id/vote", h.AuthMiddleware(h.ClearVote))
	router.GET("/api/options/:id/voters", h.AuthMiddleware(h.GetVoters))

	// Notification routes
	router.GET("/api/notifications", h.AuthMiddleware(h.GetNotifications))
	router.GET("/api/notifications/unread-count", h.AuthMiddleware(h.GetUnreadCount))
	router.POST("/api/notifications/mark-all-read", h.AuthMiddleware(h.MarkAllNotificationsRead))
	router.PUT("/api/notifications/:id/read", h.AuthMiddleware(h.MarkNotificationRead))

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL, "http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
