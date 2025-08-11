package categories

import (
    "encoding/json"
    "net/http"

    "github.com/mytheresa/go-hiring-challenge/models"
)

type CategoriesHandler struct {
    repo *models.CategoriesRepository
}

func NewCategoriesHandler(r *models.CategoriesRepository) *CategoriesHandler {
    return &CategoriesHandler{repo: r}
}

// GET /categories
func (h *CategoriesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
    categories, err := h.repo.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(categories)
}

// POST /categories
func (h *CategoriesHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
    var cat models.Category
    if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }
    if cat.Name == "" {
        http.Error(w, "El nombre es obligatorio", http.StatusBadRequest)
        return
    }
    if err := h.repo.Create(&cat); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(cat)
}