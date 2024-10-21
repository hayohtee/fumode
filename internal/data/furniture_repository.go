package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// FurnitureRepository is a type which wraps around a sql.DB connection pool
// and provide methods for creating and managing furniture to and from
// the database.
type FurnitureRepository struct {
	DB *sql.DB
}

// Insert a furniture record to the database.
func (f FurnitureRepository) Insert(furniture *Furniture) error {
	var categoryID int64
	queryCategory := `
		INSERT INTO category(name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
		RETURNING category_id`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := f.DB.QueryRowContext(ctx, queryCategory, furniture.Category).Scan(&categoryID)
	if err != nil {
		switch {
		// If no new rows was inserted (category already exist), retrieve the existing category id
		case errors.Is(err, sql.ErrNoRows):
			queryCategory = `
				SELECT category_id FROM category
				WHERE name = $1`
			err = f.DB.QueryRowContext(ctx, queryCategory, furniture.Category).Scan(&categoryID)
			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	queryFurniture := `
		INSERT INTO furniture(name, description, price, stock, banner_url, image_urls, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING furniture_id, version`

	args := []any{
		furniture.Name,
		furniture.Description,
		furniture.Price,
		furniture.Stock,
		furniture.BannerURL,
		furniture.ImageURLs,
		categoryID,
	}

	return f.DB.QueryRowContext(ctx, queryFurniture, args...).Scan(
		&furniture.FurnitureID,
		&furniture.Version,
	)
}

// GetByID retrieve a specific furniture record from the database
// given the id.
func (f FurnitureRepository) GetByID(id int64) (Furniture, error) {
	query := `
		SELECT 
			f.furniture_id, 
			f.name, 
			f.description, 
			f.price, 
			f.stock, 
			f.banner_url, 
			f.image_urls, 
			c.name AS category,
			COALESCE(AVG(r.rating), 0) AS rating,
			f.version
		FROM 
		    furniture f
		INNER JOIN 
		    category c ON f.category_id = c.category_id
		LEFT JOIN 
			review r ON f.furniture_id = r.furniture_id
		GROUP BY 
		    f.furniture_id, f.name, f.price
		WHERE f.furniture_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var furniture Furniture
	err := f.DB.QueryRowContext(ctx, query, id).Scan(
		&furniture.FurnitureID,
		&furniture.Name,
		&furniture.Description,
		&furniture.Price,
		&furniture.Stock,
		&furniture.BannerURL,
		&furniture.ImageURLs,
		&furniture.Category,
		&furniture.Rating,
		&furniture.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Furniture{}, ErrRecordNotFound
		default:
			return Furniture{}, err
		}
	}
	return furniture, nil
}

// GetAll returns list of furniture, rating and metadata
// from the database matching the provided parameters.
func (f FurnitureRepository) GetAll(name, category string, price float64, filters Filters) ([]Furniture, float32, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT 
		    COUNT(*) OVER(),
			f.furniture_id,
			f.name,
			f.price
			c.name,
			COALESCE(AVG(r.rating), 0) AS rating
		FROM
		    furniture f
		INNER JOIN 
			category c ON f.category_id = c.category_id
		LEFT JOIN 
			review r ON f.furniture_id = r.furniture_id
		GROUP BY 
		    f.furniture_id, f.name, f.price
		WHERE (to_tsvector('simple', f.name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND c.name = $2 OR $2 = ''
		AND f.price > $3
		ORDER BY %s %s, f.furniture_id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{name, category, price, filters.limit(), filters.offset()}
	rows, err := f.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, Metadata{}, err
	}
	defer rows.Close()

	var furniture []Furniture
	var totalRecords = 0
	var rating float32 = 0.0

	for rows.Next() {
		var f Furniture
		err := rows.Scan(
			&totalRecords,
			&f.FurnitureID,
			&f.Name,
			&f.Price,
			&f.Category,
			&rating,
		)
		if err != nil {
			return nil, 0, Metadata{}, err
		}
		furniture = append(furniture, f)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return furniture, rating, metadata, nil
}
