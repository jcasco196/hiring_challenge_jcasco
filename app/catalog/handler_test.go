package catalog

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gorilla/mux"
    "github.com/mytheresa/go-hiring-challenge/models"
    "github.com/shopspring/decimal"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect database: %v", err)
    }
    db.AutoMigrate(&models.Category{}, &models.Product{}, &models.Variant{})
    return db
}

func TestHandleGetByCode(t *testing.T) {
    db := setupTestDB(t)
    cat := models.Category{Code: "clothing", Name: "Clothing"}
    db.Create(&cat)
    prod := models.Product{Code: "PROD001", Price: mustDecimal("10.99"), CategoryID: cat.ID}
    db.Create(&prod)
    db.Create(&models.Variant{ProductID: prod.ID, Name: "Variant A", SKU: "SKU001A", Price: mustDecimalPtr("11.99")})
    db.Create(&models.Variant{ProductID: prod.ID, Name: "Variant B", SKU: "SKU001B", Price: nil})

    repo := models.NewProductsRepository(db)
    handler := NewCatalogHandler(repo)

    req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
    w := httptest.NewRecorder()

    // Simula mux.Vars
    req = muxSetVars(req, map[string]string{"code": "PROD001"})

    handler.HandleGetByCode(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }

    var resp ProductDetail
    if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
        t.Fatalf("error decoding response: %v", err)
    }
    if resp.Code != "PROD001" {
        t.Errorf("expected code 'PROD001', got '%s'", resp.Code)
    }
    if resp.Category != "Clothing" {
        t.Errorf("expected category 'Clothing', got '%s'", resp.Category)
    }
    if len(resp.Variants) != 2 {
        t.Errorf("expected 2 variants, got %d", len(resp.Variants))
    }
    // Verifica herencia de precio
    if resp.Variants[1].Price != 10.99 {
        t.Errorf("expected variant price 10.99, got %f", resp.Variants[1].Price)
    }
}

// Helpers para decimal
func mustDecimal(val string) decimal.Decimal {
    d, _ := decimal.NewFromString(val)
    return d
}
func mustDecimalPtr(val string) *decimal.Decimal {
    d := mustDecimal(val)
    return &d
}

// Simula mux.Vars en tests
func muxSetVars(r *http.Request, vars map[string]string) *http.Request {
    return mux.SetURLVars(r, vars)
}