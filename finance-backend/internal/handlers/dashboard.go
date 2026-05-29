package handlers

import (
    "finance-backend/internal/middleware"
    "finance-backend/internal/models"
    "finance-backend/internal/storage"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

type DashboardHandler struct {
    repo *storage.Repository
}

func NewDashboardHandler(repo *storage.Repository) *DashboardHandler {
    return &DashboardHandler{repo: repo}
}

// GetStats возвращает основную статистику для дашборда
func (h *DashboardHandler) GetStats(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Парсим месяц из запроса
    monthStr := c.Query("month")
    var year, month int
    now := time.Now()

    if monthStr != "" {
        t, err := time.Parse("2006-01", monthStr)
        if err == nil {
            year = t.Year()
            month = int(t.Month())
        } else {
            year = now.Year()
            month = int(now.Month())
        }
    } else {
        year = now.Year()
        month = int(now.Month())
    }

    // 1. Получаем общий баланс
    balance, err := h.repo.GetBalance(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
        return
    }

    // 2. Получаем доходы и расходы за месяц
    income, expense, err := h.repo.GetMonthlyStats(userID, year, month)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get monthly stats"})
        return
    }

    // 3. Получаем категории с прогрессом бюджета ← НОВОЕ
    categoryBudgets, err := h.repo.GetCategoriesWithBudget(userID, year, month)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get category budgets"})
        return
    }

    // 4. Получаем последние 5 транзакций
    recentTransactions, err := h.repo.GetTransactions(userID, nil, nil, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recent transactions"})
        return
    }

    if recentTransactions == nil {
        recentTransactions = []models.TransactionWithCategory{}
    }
    if len(recentTransactions) > 5 {
        recentTransactions = recentTransactions[:5]
    }

    // 5. Формируем ответ
    stats := models.DashboardStats{
        Balance:            balance,
        TotalIncome:        income,
        TotalExpense:       expense,
        CategoryBudgets:    categoryBudgets,      // ← НОВОЕ
        RecentTransactions: recentTransactions,
    }

    c.JSON(http.StatusOK, stats)
}

// GetExpensesByCategory возвращает расходы по категориям для круговой диаграммы
// GET /api/charts/expenses-by-category
// GetExpensesByCategory возвращает расходы по категориям для круговой диаграммы
func (h *DashboardHandler) GetExpensesByCategory(c *gin.Context) {
    userID := middleware.GetUserID(c)

    monthStr := c.Query("month")
    var year, month int
    now := time.Now()

    if monthStr != "" {
        t, err := time.Parse("2006-01", monthStr)
        if err == nil {
            year = t.Year()
            month = int(t.Month())
        } else {
            year = now.Year()
            month = int(now.Month())
        }
    } else {
        year = now.Year()
        month = int(now.Month())
    }

    expenses, err := h.repo.GetCategoryExpenses(userID, year, month)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get category expenses",
        })
        return
    }

    // Гарантируем, что возвращается массив, даже если nil
    if expenses == nil {
        expenses = []models.CategoryExpense{}
    }

    c.JSON(http.StatusOK, expenses)
}

// GetDailyTrend возвращает динамику доходов/расходов по дням
// GET /api/charts/daily-trend
func (h *DashboardHandler) GetDailyTrend(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Парсим месяц из запроса
    monthStr := c.Query("month")
    var year, month int
    now := time.Now()

    if monthStr != "" {
        t, err := time.Parse("2006-01", monthStr)
        if err == nil {
            year = t.Year()
            month = int(t.Month())
        } else {
            year = now.Year()
            month = int(now.Month())
        }
    } else {
        year = now.Year()
        month = int(now.Month())
    }

    trends, err := h.repo.GetDailyTrend(userID, year, month)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get daily trend",
        })
        return
    }

    // Гарантируем, что возвращается массив, даже если nil
    if trends == nil {
        trends = []models.DailyTrend{}
    }

    c.JSON(http.StatusOK, trends)
}
