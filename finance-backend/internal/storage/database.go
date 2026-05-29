package storage

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
    // Формируем строку подключения для PostgreSQL
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    log.Println("PostgreSQL database connected successfully")

    if err = runMigrations(); err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    return nil
}

func runMigrations() error {
    migrations := []string{
        // Таблица пользователей
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,

        // Таблица категорий
        `CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            name VARCHAR(100) NOT NULL,
            type VARCHAR(10) CHECK (type IN ('income', 'expense')) NOT NULL,
            monthly_limit DECIMAL(10,2) DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,

        // Таблица транзакций
        `CREATE TABLE IF NOT EXISTS transactions (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
            amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
            date DATE NOT NULL,
            note TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,

        // Индексы для ускорения запросов
        `CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, date)`,
        `CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_id)`,
        `CREATE INDEX IF NOT EXISTS idx_categories_user ON categories(user_id)`,
    }

    for _, migration := range migrations {
        if _, err := DB.Exec(migration); err != nil {
            log.Printf("Migration error: %v", err)
            return err
        }
    }

    log.Println("Migrations completed successfully")
    return nil
}

func CloseDB() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}
