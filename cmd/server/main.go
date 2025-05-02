package main

import (
	"database/sql"
	"github.com/W1mer/TicketControlSystem/internal/api"
	"github.com/W1mer/TicketControlSystem/internal/repository"
	"github.com/W1mer/TicketControlSystem/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация логгера
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Подключение к PostgreSQL
	connStr := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Инициализация слоев
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := api.NewTaskHandler(taskService, logger, os.Getenv("API_KEY"))

	// Настройка маршрутов
	router := mux.NewRouter()
	router.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// Запуск сервера
	logger.Info().Msg("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
