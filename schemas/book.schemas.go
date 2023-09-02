package schemas

type SearchBooksSchema struct {
	SearchBy string `json:"searchBy" validate:"required,oneof=name author"`
	Search   string `json:"search" validate:"required"`
}

type FilterBooksSchema struct {
	Categories []string  `json:"categories" validate:"omitempty,dive,max=255"`
	PriceRange []float64 `json:"priceRange" validate:"omitempty,dive,min=0"`
	InStock    bool      `json:"inStock"`
}

type CreateBookSchema struct {
	Name        string   `json:"name" validate:"required"`
	Author      string   `json:"author" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Categories  []string `json:"categories" validate:"required"`
	Cost        float64  `json:"cost" validate:"required"`
	Stock       int      `json:"stock" validate:"required"`
}

type UpdateBookSchema struct {
	Name        string   `json:"name,omitempty" validate:"omitempty,required"`
	Author      string   `json:"author,omitempty" validate:"omitempty,required"`
	Description string   `json:"description,omitempty" validate:"omitempty,required"`
	Categories  []string `json:"categories,omitempty" validate:"omitempty,required"`
	Cost        float64  `json:"cost,omitempty" validate:"omitempty,required"`
	Stock       int      `json:"stock,omitempty" validate:"omitempty,required"`
}

type DeleteBooksSchema struct {
	IDs []string `json:"ids" validate:"required"`
}
