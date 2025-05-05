package models

import "database/sql"

func (p *Product) ToCache() ProductCache {
	var catID *int64
	if p.CategoryID.Valid {
		catID = &p.CategoryID.Int64
	}
	return ProductCache{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		Price:        p.Price,
		Availability: p.Availability,
		Url:          p.Url.String,
		CategoryID:   catID,
	}
}

func (c *ProductCache) ToModel() Product {
	product := Product{
		ID:           c.ID,
		Name:         c.Name,
		Description:  c.Description,
		Price:        c.Price,
		Availability: c.Availability,
		Url:          sql.NullString{String: c.Url, Valid: c.Url != ""},
	}
	if c.CategoryID != nil {
		product.CategoryID = sql.NullInt64{Int64: *c.CategoryID, Valid: true}
	}
	return product
}

type ProductCache struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Availability bool    `json:"availability"`
	Url          string  `json:"url"`
	CategoryID   *int64  `json:"category_id"`
}

func ConvertProductsToCache(products []Product) []ProductCache {
	cached := make([]ProductCache, 0, len(products))
	for _, p := range products {
		cached = append(cached, p.ToCache())
	}
	return cached
}

func ConvertCacheToProducts(cached []ProductCache) []Product {
	products := make([]Product, 0, len(cached))
	for _, c := range cached {
		products = append(products, c.ToModel())
	}
	return products
}
