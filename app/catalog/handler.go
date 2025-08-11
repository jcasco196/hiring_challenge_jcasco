package catalog

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/gorilla/mux"
    "github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
    Products []Product `json:"products"`
}

type Product struct {
    Code     string  `json:"code"`
    Price    float64 `json:"price"`
    Category string  `json:"category"`
}

type Variant struct {
    Name  string  `json:"name"`
    SKU   string  `json:"sku"`
    Price float64 `json:"price"`
}

type ProductDetail struct {
    Code     string    `json:"code"`
    Price    float64   `json:"price"`
    Category string    `json:"category"`
    Variants []Variant `json:"variants"`
}

type CatalogHandler struct {
    repo *models.ProductsRepository
}

func NewCatalogHandler(r *models.ProductsRepository) *CatalogHandler {
    return &CatalogHandler{
        repo: r,
    }
}

// Handler para GET /catalog
func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "total":    total,
        "products": products,
    })
}

// Handler para GET /catalog/{code}
func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    code := vars["code"]
    code = strings.TrimSpace(code)

    product, err := h.repo.GetProductByCode(code)
    if err != nil {
        http.Error(w, "Producto no encontrado", http.StatusNotFound)
        return
    }

    // Mapear variantes, heredando precio si es NULL
    var variants []Variant
    for _, v := range product.Variants {
        price := product.Price.InexactFloat64()
        if v.Price.Valid {
            price = v.Price.InexactFloat64()
        }
        variants = append(variants, Variant{
            Name:  v.Name,
            SKU:   v.SKU,
            Price: price,
        })
    }

    resp := ProductDetail{
        Code:     product.Code,
        Price:    product.Price.InexactFloat64(),
        Category: product.Category.Name,
        Variants: variants,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}