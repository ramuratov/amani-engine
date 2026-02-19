package postgres

import (
	"context"
	"database/sql"
)

// Product представляет структуру товара в коде, точно такую же, как в нашей таблице БД
type Product struct {
	SKU           string
	Category      string
	SizeName      string
	BustFull      float64
	WaistFull     float64
	HipsFull      float64
	ProductLength float64
}

// ProductRepository будет управлять связью с базой
type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// SaveProduct — эта функция как раз позволит боту или тебе записывать данные в таблицу
func (r *ProductRepository) SaveProduct(ctx context.Context, p Product) error {
	query := `
		INSERT INTO product_metadata (sku, category, size_name, bust_full, waist_full, hips_full, product_length)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (sku, size_name) DO UPDATE SET
			bust_full = EXCLUDED.bust_full,
			waist_full = EXCLUDED.waist_full,
			hips_full = EXCLUDED.hips_full;`

	_, err := r.db.ExecContext(ctx, query, p.SKU, p.Category, p.SizeName, p.BustFull, p.WaistFull, p.HipsFull, p.ProductLength)
	return err
}
