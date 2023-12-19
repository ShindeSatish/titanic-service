package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shindesatish/titanic-service/internal/app/handler"
	"github.com/shindesatish/titanic-service/internal/app/repository"
	"github.com/shindesatish/titanic-service/internal/app/service"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/shindesatish/titanic-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// @title Titanic Service API
// @version 1.0
// @description API for accessing Titanic passenger data
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @host localhost:8080
// @BasePath /v1
func main() {

	loadEnv()

	// Use CSV repository or SQLite repository based on configuration
	var repo repository.Repository
	if useSQLite := os.Getenv("USE_SQLITE"); useSQLite == "true" {
		// Initialize SQLite database
		db, err := sql.Open("sqlite3", "./datastore/titanic.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		// Initialize SQLite repository
		sqliteRepo := repository.NewSQLiteRepository(db)
		repo = sqliteRepo
	} else {
		// Initialize CSV repository
		csvRepo := repository.NewCSVRepository("./datastore/titanic.csv")
		repo = csvRepo
	}

	// Initialize Passenger service
	passengerService := service.NewPassengerService(repo)

	// Initialize Gin
	router := gin.Default()

	// Initialize Passenger handler
	passengerHandler := handler.NewPassengerHandler(passengerService)

	// Register routes
	v1 := router.Group("/v1")
	{
		v1.GET("/passengers", passengerHandler.GetAllPassengersHandler)
		v1.GET("/passengers/:id", passengerHandler.GetPassengerByIDHandler)
		v1.GET("/passenger-attributes/:id", passengerHandler.GetPassengerAttributesHandler)
		// Add a new route for histogram functionality
		v1.GET("/fare-histogram", passengerHandler.GetFareHistogramHandler)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
