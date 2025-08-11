package models

import "gorm.io/gorm"

type CategoriesRepository struct {
    db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
    return &CategoriesRepository{db: db}
}

func (r *CategoriesRepository) GetAll() ([]Category, error) {
    var categories []Category
    if err := r.db.Find(&categories).Error; err != nil {
        return nil, err
    }
    return categories, nil
}

func (r *CategoriesRepository) Create(cat *Category) error {
    return r.db.Create(cat).Error
}