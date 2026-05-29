package handlers

import (
    "finance-backend/internal/auth"
    "finance-backend/internal/models"
    "finance-backend/internal/storage"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    repo *storage.Repository
}

func NewAuthHandler(repo *storage.Repository) *AuthHandler {
    return &AuthHandler{repo: repo}
}

// Register обрабатывает регистрацию нового пользователя
func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Проверяем, существует ли пользователь
    existing, _ := h.repo.GetUserByEmail(req.Email)
    if existing != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    }

    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Создаём пользователя
    userID, err := h.repo.CreateUser(req.Email, string(hashedPassword))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Создаём JWT токен
    token, err := auth.GenerateToken(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Создаём категории по умолчанию для нового пользователя
    defaultCategories := []struct {
        name  string
        typ   string
        limit float64
    }{
        {"Продукты", "expense", 15000},
        {"Транспорт", "expense", 5000},
        {"Кафе", "expense", 3000},
        {"Развлечения", "expense", 4000},
        {"Зарплата", "income", 0},
        {"Фриланс", "income", 0},
    }

    for _, cat := range defaultCategories {
        h.repo.CreateCategory(userID, cat.name, cat.typ, cat.limit)
    }

    c.JSON(http.StatusCreated, gin.H{
        "token": token,
        "user": gin.H{
            "id":    userID,
            "email": req.Email,
        },
    })
}

// Login обрабатывает вход пользователя
func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Ищем пользователя
    user, err := h.repo.GetUserByEmail(req.Email)
    if err != nil || user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Проверяем пароль
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Создаём токен
    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":    user.ID,
            "email": user.Email,
        },
    })
}
