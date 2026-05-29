package main

import (
    "finance-backend/internal/handlers"
    "finance-backend/internal/middleware"
    "finance-backend/internal/service"
    "finance-backend/internal/storage"
    "log"
    "os"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Загружаем .env файл
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Подключаемся к базе данных
    if err := storage.InitDB(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer storage.CloseDB()

    // Создаём репозиторий
    repo := storage.NewRepository()

    // Создаём BudgetChecker
    budgetChecker := service.NewBudgetChecker(repo)

    // Создаём обработчики
    authHandler := handlers.NewAuthHandler(repo)
    categoryHandler := handlers.NewCategoryHandler(repo)
    transactionHandler := handlers.NewTransactionHandler(repo, budgetChecker)
    dashboardHandler := handlers.NewDashboardHandler(repo)
    exportHandler := handlers.NewExportHandler(repo)

    // Настраиваем Gin
    r := gin.Default()

    // Настройка CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Публичные маршруты (не требуют авторизации)
    auth := r.Group("/api")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }

    // Защищённые маршруты (требуют JWT токен)
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware())
    {
        // Тестовый маршрут
        api.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "pong", "user_id": middleware.GetUserID(c)})
        })

        // Категории
        api.GET("/categories", categoryHandler.GetCategories)
        api.POST("/categories", categoryHandler.CreateCategory)
        api.PUT("/categories/:id", categoryHandler.UpdateCategory)
        api.DELETE("/categories/:id", categoryHandler.DeleteCategory)

        // Транзакции
        api.GET("/transactions", transactionHandler.GetTransactions)
        api.POST("/transactions", transactionHandler.CreateTransaction)
        api.PUT("/transactions/:id", transactionHandler.UpdateTransaction)
        api.DELETE("/transactions/:id", transactionHandler.DeleteTransaction)

        // Дашборд и статистика
        api.GET("/dashboard/stats", dashboardHandler.GetStats)
        api.GET("/charts/expenses-by-category", dashboardHandler.GetExpensesByCategory)
        api.GET("/charts/daily-trend", dashboardHandler.GetDailyTrend)
        api.GET("/export/csv", exportHandler.ExportCSV)
    }

    // Запускаем сервер
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
