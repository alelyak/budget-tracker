package service

import (
    "finance-backend/internal/models"
    "finance-backend/internal/storage"
    "time"
)

type BudgetChecker struct {
    repo *storage.Repository
}

func NewBudgetChecker(repo *storage.Repository) *BudgetChecker {
    return &BudgetChecker{repo: repo}
}

// CheckBudget проверяет, не превысит ли новая транзакция лимит категории
func (bc *BudgetChecker) CheckBudget(userID, categoryID int64, amount float64, date time.Time) (*models.BudgetCheckResult, error) {
    // Получаем категорию
    category, err := bc.repo.GetCategoryByID(categoryID, userID)
    if err != nil {
        return nil, err
    }
    if category == nil {
        return nil, nil
    }

    // Проверяем только категории расходов с лимитом > 0
    if category.Type != "expense" || category.MonthlyLimit <= 0 {
        return nil, nil
    }

    // Получаем текущие расходы за месяц по этой категории
    year, month, _ := date.Date()
    spentMap, err := bc.repo.GetMonthlySpentByCategory(userID, year, int(month))
    if err != nil {
        return nil, err
    }

    currentSpent := spentMap[categoryID]
    newTotal := currentSpent + amount

    // Проверяем превышение
    if newTotal > category.MonthlyLimit {
        return &models.BudgetCheckResult{
            IsExceeded:   true,
            CategoryID:   category.ID,
            CategoryName: category.Name,
            Spent:        currentSpent,
            Limit:        category.MonthlyLimit,
            ExceedAmount: newTotal - category.MonthlyLimit,
        }, nil
    }

    return &models.BudgetCheckResult{
        IsExceeded:   false,
        CategoryID:   category.ID,
        CategoryName: category.Name,
        Spent:        currentSpent,
        Limit:        category.MonthlyLimit,
        ExceedAmount: 0,
    }, nil
}
