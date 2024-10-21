package main

import (
	"context"
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/validator"
	"net/http"
	"strconv"
	"time"
)

func (app *application) createFurnitureHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form with a 10 MB max memory limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, "must be a multipart form")
		return
	}

	name := r.Form.Get("name")
	description := r.Form.Get("description")
	priceStr := r.Form.Get("price")
	stockStr := r.Form.Get("stock")
	category := r.Form.Get("category")

	v := validator.New()
	v.Check(name != "", "name", "must be provided")
	v.Check(description != "", "description", "must be provided")
	v.Check(category != "", "category", "must be provided")
	v.Check(priceStr != "", "price", "must be provided")
	v.Check(stockStr != "", "stock", "must be provided")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		v.AddError("price", "must be a valid number")
	}

	stock, err := strconv.ParseInt(stockStr, 10, 32)
	if err != nil {
		v.AddError("stock", "must be a valid number")
	}

	_, bannerHeader, err := r.FormFile("banner")
	if err != nil {
		v.AddError("banner", "must be a valid file")
	}

	images := r.MultipartForm.File["images"]
	if images == nil || len(images) == 0 {
		v.AddError("images", "must contain at least one image")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	furniture := data.Furniture{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       int(stock),
		Category:    category,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	bannerUrl, err := app.s3Uploader.UploadImage(ctx, bannerHeader)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	imageUrls, err := app.s3Uploader.UploadImages(ctx, images)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	furniture.BannerURL = bannerUrl
	furniture.ImageURLs = imageUrls

	err = app.repositories.Furniture.Insert(&furniture)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"furniture": furniture}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllFurnitureHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Category string
		Price    float64
		Rating   float32
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Category = app.readString(qs, "category", "")
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "name")
	input.SortSafeList = []string{"name", "price", "-name", "-price"}

	priceInput, err := app.readDecimal(qs, "price", 0)
	if err != nil {
		v.AddError("price", "must be a valid decimal number")
	}
	if priceInput < 0 {
		v.AddError("price", "must be a positive number")
	}
	input.Price = priceInput

	ratingInput, err := app.readDecimal(qs, "rating", 0)
	if err != nil {
		v.AddError("rating", "must be a valid decimal number")
	}
	if ratingInput < 0 || ratingInput > 5 {
		v.AddError("rating", "must be between 0 and 5")
	}
	input.Rating = float32(ratingInput)

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	furniture, rating, metadata, err := app.repositories.Furniture.GetAll(
		input.Name,
		input.Category,
		input.Price,
		input.Filters,
	)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	items := make([]furnitureItem, len(furniture))

	for i, f := range furniture {
		items[i] = furnitureItem{
			ID:       f.FurnitureID,
			Name:     f.Name,
			Price:    f.Price,
			Rating:   rating,
			Category: f.Category,
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "items": items}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

type furnitureItem struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Rating   float32 `json:"rating"`
}
