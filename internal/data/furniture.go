package data

// Furniture is a struct that holds information about
// a specific furniture.
type Furniture struct {
	FurnitureID int64
	Name        string
	Description string
	Price       float64
	Stock       int
	BannerURL   string
	ImageURLs   []string
	Category    string
	Rating      float32
	Version     int
}
