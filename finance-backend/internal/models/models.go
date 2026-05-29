package models

import "time"

// Модель пользователя
type User struct {
    ID           int64     `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"` // "-" означает, что поле не будет в JSON
    CreatedAt    time.Time `json:"created_at"`
}

// Модель категории
type Category struct {
    ID           int64     `json:"id"`
    UserID       int64     `json:"user_id"`
    Name         string    `json:"name"`
    Type         string    `json:"type"` // "income" или "expense"
    MonthlyLimit float64   `json:"monthly_limit"`
    CreatedAt    time.Time `json:"created_at"`
}

// Модель транзакции
type Transaction struct {
    ID         int64     `json:"id"`
    UserID     int64     `json:"user_id"`
    CategoryID int64     `json:"category_id"`
    Amount     float64   `json:"amount"`
    Date       time.Time `json:"date"`
    Note       string    `json:"note"`
    CreatedAt  time.Time `json:"created_at"`
}

// Расширенная категория с потраченной суммой (для дашборда)
type CategoryWithSpent struct {
    Category
    Spent      float64 `json:"spent"`
    Remaining  float64 `json:"remaining"`
    Percent    float64 `json:"percent"`
    IsExceeded bool    `json:"is_exceeded"`
}

// Статистика для дашборда
type DashboardStats struct {
    Balance            float64                `json:"balance"`
    TotalIncome        float64                `json:"total_income"`
    TotalExpense       float64                `json:"total_expense"`
    CategoryBudgets    []CategoryWithSpent    `json:"category_budgets"`
    RecentTransactions []TransactionWithCategory `json:"recent_transactions"`
}

// Транзакция с именем категории
type TransactionWithCategory struct {
    Transaction
    CategoryName string `json:"category_name"`
    CategoryType string `json:"category_type"`
}

// Запросы от клиента
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type CategoryRequest struct {
    Name         string  `json:"name" binding:"required"`
    Type         string  `json:"type" binding:"required,oneof=income expense"`
    MonthlyLimit float64 `json:"monthly_limit"`
}

type TransactionRequest struct {
    CategoryID int64   `json:"category_id" binding:"required"`
    Amount     float64 `json:"amount" binding:"required,gt=0"`
    Date       string  `json:"date" binding:"required"`
    Note       string  `json:"note"`
}

// Результат проверки бюджета
type BudgetCheckResult struct {
    IsExceeded   bool    `json:"is_exceeded"`
    CategoryID   int64   `json:"category_id"`
    CategoryName string  `json:"category_name"`
    Spent        float64 `json:"spent"`
    Limit        float64 `json:"limit"`
    ExceedAmount float64 `json:"exceed_amount"`
}

// CategoryExpense - расходы по категории (для круговой диаграммы)
type CategoryExpense struct {
    CategoryName string  `json:"category_name"`
    Amount       float64 `json:"amount"`
    Percentage   float64 `json:"percentage"`
}

// DailyTrend - данные для графика динамики по дням
type DailyTrend struct {
    Date    string  `json:"date"`
    Income  float64 `json:"income"`
    Expense float64 `json:"expense"`
}
