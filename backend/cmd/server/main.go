package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/luisteixeira/waypoint/backend/internal/handler"
	wmiddleware "github.com/luisteixeira/waypoint/backend/internal/middleware"
	"github.com/luisteixeira/waypoint/backend/internal/repository/postgres"
	"github.com/luisteixeira/waypoint/backend/internal/service"
)

func main() {
	db := initDB()
	defer db.Close()

	activityRepo := postgres.NewPostgresActivityRepo(db)
	defRepo := postgres.NewPostgresDefinitionRepo(db)

	activityService := service.NewActivityService(activityRepo, defRepo)
	activityHandler := handler.NewActivityHandler(activityService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Waypoint Backend: Active"))
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(wmiddleware.TenantMiddleware)
		r.Post("/activities/plan", activityHandler.PlanActivity)
	})

	port := ":8080"
	log.Printf("Starting server on %s...", port)

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error parsing connection stirng: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
