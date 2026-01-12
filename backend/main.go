package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"poll_app/ent"
	"poll_app/handlers"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

func main() {
	// Initialize database connection
	db, err := sql.Open("sqlite", "file:poll_app.db?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("failed enabling foreign keys: %v", err)
	}

	// Create an ent driver from the sql.DB
	drv := entsql.OpenDB(dialect.SQLite, db)
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
	router.GET("/api/options/:id/voters", h.AuthMiddleware(h.GetVoters))

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
