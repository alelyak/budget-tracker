package handlers

import (
    "encoding/csv"
    "finance-backend/internal/middleware"
    "finance-backend/internal/storage"
    "net/http"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
)

type ExportHandler struct {
    repo *storage.Repository
}

func NewExportHandler(repo *storage.Repository) *ExportHandler {
    return &ExportHandler{repo: repo}
}

// ExportCSV экспортирует транзакции в CSV файл
// GET /api/export/csv
func (h *ExportHandler) ExportCSV(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Парсим параметры фильтрации
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

    // Получаем транзакции с фильтрами
    transactions, err := h.repo.GetTransactions(userID, startDate, endDate, categoryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch transactions",
        })
        return
    }

    // Устанавливаем заголовки для CSV файла
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", "attachment;filename=transactions_"+strconv.FormatInt(time.Now().Unix(), 10)+".csv")

    // Создаём CSV writer
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()

    // Записываем заголовки
    headers := []string{"Date", "Category", "Type", "Amount", "Note"}
    if err := writer.Write(headers); err != nil {
        return
    }

    // Записываем данные
    for _, t := range transactions {
        // Определяем знак суммы
        amountStr := strconv.FormatFloat(t.Amount, 'f', 2, 64)
        if t.CategoryType == "expense" {
            amountStr = "-" + amountStr
        } else {
            amountStr = "+" + amountStr
        }

        row := []string{
            t.Date.Format("2006-01-02"),
            t.CategoryName,
            t.CategoryType,
            amountStr,
            t.Note,
        }
        if err := writer.Write(row); err != nil {
            return
        }
    }
}
