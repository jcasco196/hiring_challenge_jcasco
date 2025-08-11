package catalog

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
    Products []Product `json:"products"`
}

type Product struct {
    Code     string  `json:"code"`
    Price    float64 `json:"price"`
    Category string  `json:"category"` // Nuevo campo para la categoría
}

type CatalogHandler struct {
    repo *models.ProductsRepository
}

func NewCatalogHandler(r *models.ProductsRepository) *CatalogHandler {
    return &CatalogHandler{
        repo: r,
    }
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
    // Leer parámetros de paginación y filtros
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    category := r.URL.Query().Get("category")
    priceLtStr := r.URL.Query().Get("price_lt")
    var priceLt float64
    if priceLtStr != "" {
        priceLt, _ = strconv.ParseFloat(priceLtStr, 64)
    }

    res, total, err := h.repo.GetProductsFiltered(offset, limit, category, priceLt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Map response
    products := make([]Product, len(res))
    for i, p := range res {
        categoryName := ""
        if p.Category.Name != "" {
            categoryName = p.Category.Name
        }
        products[i] = Product{
            Code:     p.Code,
            Price:    p.Price.InexactFloat64(),
            Category: categoryName,
        }
    }

    // Respuesta con total y productos
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "total":    total,
        "products": products,
    })
}