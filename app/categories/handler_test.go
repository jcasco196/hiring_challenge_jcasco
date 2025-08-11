package categories

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/mytheresa/go-hiring-challenge/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect database: %v", err)
    }
    db.AutoMigrate(&models.Category{})
    return db
}

func TestHandleGetCategories(t *testing.T) {
    db := setupTestDB(t)
    repo := models.NewCategoriesRepository(db)
    handler := NewCategoriesHandler(repo)

    // Insert test data
    db.Create(&models.Category{Code: "cat1", Name: "TestCat1"})
    db.Create(&models.Category{Code: "cat2", Name: "TestCat2"})

    req := httptest.NewRequest("GET", "/categories", nil)
    w := httptest.NewRecorder()

    handler.HandleGet(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }

    var cats []models.Category
    if err := json.NewDecoder(w.Body).Decode(&cats); err != nil {
        t.Fatalf("error decoding response: %v", err)
    }
    if len(cats) != 2 {
        t.Fatalf("expected 2 categories, got %d", len(cats))
    }
}

func TestHandlePostCategories(t *testing.T) {
    db := setupTestDB(t)
    repo := models.NewCategoriesRepository(db)
    handler := NewCategoriesHandler(repo)

    cat := models.Category{Name: "NuevaCat"}
    body, _ := json.Marshal(cat)
    req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
    w := httptest.NewRecorder()

    handler.HandlePost(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected status 201, got %d", w.Code)
    }

    var resp models.Category
    if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
        t.Fatalf("error decoding response: %v", err)
    }
    if resp.Name != "NuevaCat" {
        t.Fatalf("expected category name 'NuevaCat', got '%s'", resp.Name)
    }
}