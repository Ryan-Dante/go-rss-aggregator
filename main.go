package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ryan-Dante/go-rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	// Import the postgres driver
	_ "github.com/lib/pq"
)

type apiConfig struct {
	// Holds the database connection
	DB *database.Queries
}

func main() {
	// Load environment variables from .env file
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbURL is not found in the environment variables")
	}

	// Connect to the database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to the database")
	}

	// Convert database connection
	apiCfg := &apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	// Add routes
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))

	router.Mount("/v1", v1Router)

	svr := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Starting server on port %s", portString)
	err = svr.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting server on port " + portString)

}
