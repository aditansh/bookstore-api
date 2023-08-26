package schemas

type SearchBooksSchema struct {
	SearchBy string `json:"searchBy" validate:"required; oneof=name author"`
	Search   string `json:"search" validate:"required"`
}

type FilterBooksSchema struct {
	Filters BookFiltersSchema `json:"filters" validate:"required"`
}

type BookFiltersSchema struct {
	Authors     []string  `json:"authors"`
	Categories  []string  `json:"categories"`
	PriceRange  []float64 `json:"priceRange"`
	InStock     bool      `json:"inStock"`
	RatingRange []int     `json:"ratingRange"`
	Vendors     []string  `json:"vendors"`
}

type CreateBookSchema struct {
	Name        string   `json:"name" validate:"required"`
	Author      string   `json:"author" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Categories  []string `json:"categories" validate:"required"`
	Price       float64  `json:"price" validate:"required"`
	Stock       int      `json:"stock" validate:"required"`
}

type UpdateBookSchema struct {
	ID          string   `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Author      string   `json:"author" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Categories  []string `json:"categories" validate:"required"`
	Price       float64  `json:"price" validate:"required"`
	Stock       int      `json:"stock" validate:"required"`
}

type DeleteBooksSchema struct {
	IDs []string `json:"ids" validate:"required"`
}
