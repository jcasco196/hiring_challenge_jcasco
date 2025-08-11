package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
    "github.com/mytheresa/go-hiring-challenge/app/catalog"
    "github.com/mytheresa/go-hiring-challenge/app/database"
    "github.com/mytheresa/go-hiring-challenge/models"
    "github.com/mytheresa/go-hiring-challenge/app/categories"
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(".env"); err != nil {
        log.Fatalf("Error loading .env file: %s", err)
    }

    // signal handling for graceful shutdown
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    // Initialize database connection
    db, close := database.New(
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_DB"),
        os.Getenv("POSTGRES_PORT"),
    )
    defer close()

    // Initialize handlers
    prodRepo := models.NewProductsRepository(db)
    cat := catalog.NewCatalogHandler(prodRepo)

    // Set up routing con gorilla/mux
    router := mux.NewRouter()
    router.HandleFunc("/catalog", cat.HandleGet).Methods("GET")
    router.HandleFunc("/catalog/{code}", cat.HandleGetByCode).Methods("GET")

    catRepo := models.NewCategoriesRepository(db)
    catHandler := categories.NewCategoriesHandler(catRepo)

    router.HandleFunc("/categories", catHandler.HandleGet).Methods("GET")
    router.HandleFunc("/categories", catHandler.HandlePost).Methods("POST")
    
    // Set up the HTTP server
    srv := &http.Server{
        Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
        Handler: router,
    }

    // Start the server
    go func() {
        log.Printf("Starting server on http://%s", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %s", err)
        }

        log.Println("Server stopped gracefully")
    }()

    <-ctx.Done()
    log.Println("Shutting down server...")
    srv.Shutdown(ctx)
    stop()
}