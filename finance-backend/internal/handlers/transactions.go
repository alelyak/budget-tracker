package handlers

import (
    "finance-backend/internal/middleware"
    "finance-backend/internal/models"
    "finance-backend/internal/storage"
	"finance-backend/internal/service"
	"fmt"
    "net/http"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
)

type TransactionHandler struct {
    repo *storage.Repository
	budgetChecker *service.BudgetChecker
}

func NewTransactionHandler(repo *storage.Repository, budgetChecker *service.BudgetChecker) *TransactionHandler {
    return &TransactionHandler{
        repo:          repo,
        budgetChecker: budgetChecker,
    }
}

// GetTransactions возвращает все транзакции пользователя с фильтрацией
// GET /api/transactions
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Парсим query параметры для фильтрации
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    categoryIDStr := c.Query("category_id")

    var startDate, endDate *time.Time
    if startDateStr != "" {
        t, err := time.Parse("2006-01-02", startDateStr)
        if err == nil {
            startDate = &t
        }
    }
    if endDateStr != "" {
        t, err := time.Parse("2006-01-02", endDateStr)
        if err == nil {
            endDate = &t
        }
    }

    var categoryID *int64
    if categoryIDStr != "" {
        id, err := strconv.ParseInt(categoryIDStr, 10, 64)
        if err == nil {
            categoryID = &id
        }
    }

    transactions, err := h.repo.GetTransactions(userID, startDate, endDate, categoryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch transactions",
        })
        return
    }

    c.JSON(http.StatusOK, transactions)
}

// CreateTransaction создаёт новую транзакцию
// POST /api/transactions
// CreateTransaction создаёт новую транзакцию
// POST /api/transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    userID := middleware.GetUserID(c)

    var req models.TransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Валидация
    if req.CategoryID <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Valid category ID is required",
        })
        return
    }

    if req.Amount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Amount must be greater than 0",
        })
        return
    }

    // Парсим дату
    date, err := time.Parse("2006-01-02", req.Date)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid date format, use YYYY-MM-DD",
        })
        return
    }

    // Проверяем, что категория существует и принадлежит пользователю
    category, err := h.repo.GetCategoryByID(req.CategoryID, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to verify category",
        })
        return
    }
    if category == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Category not found or does not belong to you",
        })
        return
    }

    // ★★★★★ ПРОВЕРКА БЮДЖЕТА ★★★★★
    budgetCheck, err := h.budgetChecker.CheckBudget(userID, req.CategoryID, req.Amount, date)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to check budget",
        })
        return
    }

    // Создаём транзакцию
    transactionID, err := h.repo.CreateTransaction(userID, req.CategoryID, req.Amount, date, req.Note)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create transaction",
        })
        return
    }

    // Формируем ответ
    response := gin.H{
        "id":      transactionID,
        "message": "Transaction created successfully",
    }

    // Добавляем предупреждение о бюджете, если нужно
    if budgetCheck != nil && budgetCheck.IsExceeded {
        response["budget_warning"] = gin.H{
            "category_name": budgetCheck.CategoryName,
            "limit":         budgetCheck.Limit,
            "spent":         budgetCheck.Spent + req.Amount,
            "exceed_amount": budgetCheck.ExceedAmount,
            "message":       "Превышен бюджет на категорию \"" + budgetCheck.CategoryName + "\" на " + formatMoney(budgetCheck.ExceedAmount),
        }
    }

    c.JSON(http.StatusCreated, response)
}

// formatMoney - вспомогательная функция для форматирования денег
func formatMoney(amount float64) string {
    return fmt.Sprintf("%.2f", amount)
}

// UpdateTransaction обновляет существующую транзакцию
// PUT /api/transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Получаем ID транзакции из URL
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid transaction ID",
        })
        return
    }

    // Проверяем, существует ли транзакция и принадлежит ли пользователю
    existing, err := h.repo.GetTransactionByID(id, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch transaction",
        })
        return
    }
    if existing == nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Transaction not found",
        })
        return
    }

    // Читаем данные из запроса
    var req models.TransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Валидация
    if req.CategoryID <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Valid category ID is required",
        })
        return
    }

    if req.Amount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Amount must be greater than 0",
        })
        return
    }

    // Парсим дату
    date, err := time.Parse("2006-01-02", req.Date)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid date format, use YYYY-MM-DD",
        })
        return
    }

    // Проверяем, что категория существует и принадлежит пользователю
    category, err := h.repo.GetCategoryByID(req.CategoryID, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to verify category",
        })
        return
    }
    if category == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Category not found or does not belong to you",
        })
        return
    }

    // Обновляем транзакцию
    err = h.repo.UpdateTransaction(id, userID, req.CategoryID, req.Amount, date, req.Note)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to update transaction",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Transaction updated successfully",
    })
}

// DeleteTransaction удаляет транзакцию
// DELETE /api/transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Получаем ID транзакции из URL
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid transaction ID",
        })
        return
    }

    // Проверяем, существует ли транзакция и принадлежит ли пользователю
    existing, err := h.repo.GetTransactionByID(id, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch transaction",
        })
        return
    }
    if existing == nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Transaction not found",
        })
        return
    }

    // Удаляем транзакцию
    err = h.repo.DeleteTransaction(id, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to delete transaction",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Transaction deleted successfully",
    })
}
