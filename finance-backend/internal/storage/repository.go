package storage

import (
    "database/sql"
    "finance-backend/internal/models"
    "fmt"
    "time"
)

type Repository struct{}

func NewRepository() *Repository {
    return &Repository{}
}

// ========== РАБОТА С ПОЛЬЗОВАТЕЛЯМИ ==========

func (r *Repository) CreateUser(email, passwordHash string) (int64, error) {
    var id int64
    err := DB.QueryRow(
        "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
        email, passwordHash,
    ).Scan(&id)
    return id, err
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := DB.QueryRow(
        "SELECT id, email, password_hash, created_at FROM users WHERE email = $1",
        email,
    ).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *Repository) GetUserByID(id int64) (*models.User, error) {
    var user models.User
    err := DB.QueryRow(
        "SELECT id, email, password_hash, created_at FROM users WHERE id = $1",
        id,
    ).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ========== РАБОТА С КАТЕГОРИЯМИ ==========

func (r *Repository) CreateCategory(userID int64, name, categoryType string, monthlyLimit float64) (int64, error) {
    var id int64
    err := DB.QueryRow(
        "INSERT INTO categories (user_id, name, type, monthly_limit) VALUES ($1, $2, $3, $4) RETURNING id",
        userID, name, categoryType, monthlyLimit,
    ).Scan(&id)
    return id, err
}

func (r *Repository) GetCategories(userID int64) ([]models.Category, error) {
    rows, err := DB.Query(
        "SELECT id, user_id, name, type, monthly_limit, created_at FROM categories WHERE user_id = $1 ORDER BY type, name",
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var cat models.Category
        err := rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Type, &cat.MonthlyLimit, &cat.CreatedAt)
        if err != nil {
            return nil, err
        }
        categories = append(categories, cat)
    }
    return categories, nil
}

func (r *Repository) GetCategoryByID(id, userID int64) (*models.Category, error) {
    var cat models.Category
    err := DB.QueryRow(
        "SELECT id, user_id, name, type, monthly_limit, created_at FROM categories WHERE id = $1 AND user_id = $2",
        id, userID,
    ).Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Type, &cat.MonthlyLimit, &cat.CreatedAt)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &cat, nil
}

func (r *Repository) UpdateCategory(id, userID int64, name, categoryType string, monthlyLimit float64) error {
    result, err := DB.Exec(
        "UPDATE categories SET name = $1, type = $2, monthly_limit = $3 WHERE id = $4 AND user_id = $5",
        name, categoryType, monthlyLimit, id, userID,
    )
    if err != nil {
        return err
    }

    affected, _ := result.RowsAffected()
    if affected == 0 {
        return fmt.Errorf("category not found")
    }
    return nil
}

func (r *Repository) DeleteCategory(id, userID int64) error {
    result, err := DB.Exec("DELETE FROM categories WHERE id = $1 AND user_id = $2", id, userID)
    if err != nil {
        return err
    }

    affected, _ := result.RowsAffected()
    if affected == 0 {
        return fmt.Errorf("category not found")
    }
    return nil
}

// ========== РАБОТА С ТРАНЗАКЦИЯМИ ==========

func (r *Repository) CreateTransaction(userID, categoryID int64, amount float64, date time.Time, note string) (int64, error) {
    var id int64
    err := DB.QueryRow(
        "INSERT INTO transactions (user_id, category_id, amount, date, note) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        userID, categoryID, amount, date, note,
    ).Scan(&id)
    return id, err
}

func (r *Repository) GetTransactions(userID int64, startDate, endDate *time.Time, categoryID *int64) ([]models.TransactionWithCategory, error) {
    query := `SELECT t.id, t.user_id, t.category_id, t.amount, t.date, t.note, t.created_at,
              c.name as category_name, c.type as category_type
              FROM transactions t
              JOIN categories c ON t.category_id = c.id
              WHERE t.user_id = $1`

    args := []interface{}{userID}
    argIndex := 2

    if startDate != nil {
        query += fmt.Sprintf(" AND t.date >= $%d", argIndex)
        args = append(args, *startDate)
        argIndex++
    }
    if endDate != nil {
        query += fmt.Sprintf(" AND t.date <= $%d", argIndex)
        args = append(args, *endDate)
        argIndex++
    }
    if categoryID != nil {
        query += fmt.Sprintf(" AND t.category_id = $%d", argIndex)
        args = append(args, *categoryID)
        argIndex++
    }

    query += " ORDER BY t.date DESC, t.created_at DESC LIMIT 100"

    rows, err := DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.TransactionWithCategory
    for rows.Next() {
        var t models.TransactionWithCategory
        err := rows.Scan(&t.ID, &t.UserID, &t.CategoryID, &t.Amount, &t.Date, &t.Note, &t.CreatedAt, &t.CategoryName, &t.CategoryType)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

// GetTransactionByID возвращает транзакцию по ID (с проверкой принадлежности пользователю)
func (r *Repository) GetTransactionByID(id, userID int64) (*models.TransactionWithCategory, error) {
    var t models.TransactionWithCategory
    err := DB.QueryRow(`
        SELECT t.id, t.user_id, t.category_id, t.amount, t.date, t.note, t.created_at,
               c.name as category_name, c.type as category_type
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.id = $1 AND t.user_id = $2
    `, id, userID).Scan(
        &t.ID, &t.UserID, &t.CategoryID, &t.Amount, &t.Date, &t.Note, &t.CreatedAt,
        &t.CategoryName, &t.CategoryType,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &t, nil
}

func (r *Repository) UpdateTransaction(id, userID, categoryID int64, amount float64, date time.Time, note string) error {
    result, err := DB.Exec(
        "UPDATE transactions SET category_id = $1, amount = $2, date = $3, note = $4 WHERE id = $5 AND user_id = $6",
        categoryID, amount, date, note, id, userID,
    )
    if err != nil {
        return err
    }

    affected, _ := result.RowsAffected()
    if affected == 0 {
        return fmt.Errorf("transaction not found")
    }
    return nil
}

func (r *Repository) DeleteTransaction(id, userID int64) error {
    result, err := DB.Exec("DELETE FROM transactions WHERE id = $1 AND user_id = $2", id, userID)
    if err != nil {
        return err
    }

    affected, _ := result.RowsAffected()
    if affected == 0 {
        return fmt.Errorf("transaction not found")
    }
    return nil
}

// ========== СТАТИСТИКА И АНАЛИТИКА ==========

func (r *Repository) GetMonthlySpentByCategory(userID int64, year, month int) (map[int64]float64, error) {
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, -1)

    rows, err := DB.Query(`
        SELECT category_id, COALESCE(SUM(amount), 0) as spent
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.user_id = $1 AND t.date BETWEEN $2 AND $3 AND c.type = 'expense'
        GROUP BY category_id
    `, userID, startDate, endDate)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    spent := make(map[int64]float64)
    for rows.Next() {
        var categoryID int64
        var amount float64
        if err := rows.Scan(&categoryID, &amount); err != nil {
            return nil, err
        }
        spent[categoryID] = amount
    }
    return spent, nil
}

func (r *Repository) GetBalance(userID int64) (float64, error) {
    var balance float64
    err := DB.QueryRow(`
        SELECT COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE -t.amount END), 0)
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.user_id = $1
    `, userID).Scan(&balance)
    return balance, err
}

func (r *Repository) GetMonthlyStats(userID int64, year, month int) (income, expense float64, err error) {
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, -1)

    err = DB.QueryRow(`
        SELECT
            COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE 0 END), 0) as income,
            COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount ELSE 0 END), 0) as expense
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.user_id = $1 AND t.date BETWEEN $2 AND $3
    `, userID, startDate, endDate).Scan(&income, &expense)

    return
}

// GetCategoryExpenses возвращает расходы по категориям за месяц
func (r *Repository) GetCategoryExpenses(userID int64, year, month int) ([]models.CategoryExpense, error) {
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, -1)

    rows, err := DB.Query(`
        SELECT c.name, COALESCE(SUM(t.amount), 0) as amount
        FROM categories c
        LEFT JOIN transactions t ON t.category_id = c.id AND t.date BETWEEN $1 AND $2
        WHERE c.user_id = $3 AND c.type = 'expense'
        GROUP BY c.id, c.name
        ORDER BY amount DESC
    `, startDate, endDate, userID)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var expenses []models.CategoryExpense
    var total float64

    // Сначала собираем все расходы
    for rows.Next() {
        var e models.CategoryExpense
        if err := rows.Scan(&e.CategoryName, &e.Amount); err != nil {
            return nil, err
        }
        expenses = append(expenses, e)
        total += e.Amount
    }

    // Вычисляем проценты
    for i := range expenses {
        if total > 0 {
            expenses[i].Percentage = (expenses[i].Amount / total) * 100
        }
    }

    return expenses, nil
}

// GetDailyTrend возвращает динамику доходов/расходов по дням за месяц
func (r *Repository) GetDailyTrend(userID int64, year, month int) ([]models.DailyTrend, error) {
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, -1)

    rows, err := DB.Query(`
        SELECT
            DATE(t.date) as day,
            COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE 0 END), 0) as income,
            COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount ELSE 0 END), 0) as expense
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.user_id = $1 AND t.date BETWEEN $2 AND $3
        GROUP BY DATE(t.date)
        ORDER BY day ASC
    `, userID, startDate, endDate)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var trends []models.DailyTrend
    for rows.Next() {
        var d models.DailyTrend
        var date time.Time
        if err := rows.Scan(&date, &d.Income, &d.Expense); err != nil {
            return nil, err
        }
        d.Date = date.Format("2006-01-02")
        trends = append(trends, d)
    }

    return trends, nil
}

// GetCategoriesWithBudget возвращает категории с информацией о бюджете
func (r *Repository) GetCategoriesWithBudget(userID int64, year, month int) ([]models.CategoryWithSpent, error) {
    // Получаем все категории пользователя
    categories, err := r.GetCategories(userID)
    if err != nil {
        return nil, err
    }

    // Получаем расходы за месяц по категориям
    spentMap, err := r.GetMonthlySpentByCategory(userID, year, month)
    if err != nil {
        return nil, err
    }

    // Собираем результат только для расходных категорий с лимитом
    var result []models.CategoryWithSpent
    for _, cat := range categories {
        if cat.Type == "expense" && cat.MonthlyLimit > 0 {
            spent := spentMap[cat.ID]
            remaining := cat.MonthlyLimit - spent
            percent := (spent / cat.MonthlyLimit) * 100
            if percent > 100 {
                percent = 100
            }

            result = append(result, models.CategoryWithSpent{
                Category:   cat,
                Spent:      spent,
                Remaining:  remaining,
                Percent:    percent,
                IsExceeded: spent > cat.MonthlyLimit,
            })
        }
    }

    return result, nil
}
