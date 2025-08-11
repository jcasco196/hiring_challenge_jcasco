package models

import (
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetProductsFiltered(offset, limit int, category string, priceLt float64) ([]Product, int64, error) {
    var products []Product
    var total int64

    query := r.db.Model(&Product{}).Preload("Category")

    if category != "" {
        query = query.Joins("JOIN categories ON categories.id = products.category_id").Where("categories.name = ?", category)
    }
    if priceLt > 0 {
        query = query.Where("price < ?", priceLt)
    }

    // Contar total antes de paginar
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset > 0 {
        query = query.Offset(offset)
    }

    if err := query.Find(&products).Error; err != nil {
        return nil, 0, err
    }
    return products, total, nil
}