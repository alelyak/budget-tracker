package handlers

import (
    "finance-backend/internal/middleware"
    "finance-backend/internal/models"
    "finance-backend/internal/storage"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

type CategoryHandler struct {
    repo *storage.Repository
}

func NewCategoryHandler(repo *storage.Repository) *CategoryHandler {
    return &CategoryHandler{repo: repo}
}

// GetCategories возвращает все категории пользователя
// GET /api/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
    userID := middleware.GetUserID(c)

    categories, err := h.repo.GetCategories(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch categories",
        })
        return
    }

    c.JSON(http.StatusOK, categories)
}

// CreateCategory создаёт новую категорию
// POST /api/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
    userID := middleware.GetUserID(c)

    var req models.CategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Валидация
    if req.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Category name is required",
        })
        return
    }

    if req.Type != "income" && req.Type != "expense" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Type must be 'income' or 'expense'",
        })
        return
    }

    if req.MonthlyLimit < 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Monthly limit cannot be negative",
        })
        return
    }

    categoryID, err := h.repo.CreateCategory(userID, req.Name, req.Type, req.MonthlyLimit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create category",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":      categoryID,
        "message": "Category created successfully",
    })
}

// UpdateCategory обновляет существующую категорию
// PUT /api/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Получаем ID категории из URL
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid category ID",
        })
        return
    }

    // Проверяем, существует ли категория и принадлежит ли пользователю
    existing, err := h.repo.GetCategoryByID(id, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch category",
        })
        return
    }
    if existing == nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Category not found",
        })
        return
    }

    // Читаем данные из запроса
    var req models.CategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Валидация
    if req.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Category name is required",
        })
        return
    }

    if req.Type != "income" && req.Type != "expense" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Type must be 'income' or 'expense'",
        })
        return
    }

    if req.MonthlyLimit < 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Monthly limit cannot be negative",
        })
        return
    }

    // Обновляем категорию
    err = h.repo.UpdateCategory(id, userID, req.Name, req.Type, req.MonthlyLimit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to update category",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Category updated successfully",
    })
}

// DeleteCategory удаляет категорию
// DELETE /api/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
    userID := middleware.GetUserID(c)

    // Получаем ID категории из URL
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid category ID",
        })
        return
    }

    // Проверяем, существует ли категория и принадлежит ли пользователю
    existing, err := h.repo.GetCategoryByID(id, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch category",
        })
        return
    }
    if existing == nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Category not found",
        })
        return
    }

    // Удаляем категорию
    err = h.repo.DeleteCategory(id, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to delete category",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Category deleted successfully",
    })
}
