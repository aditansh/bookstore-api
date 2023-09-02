package schemas

type SearchOrdersSchema struct {
	SearchBy string `json:"searchBy" validate:"required; oneof=bookName author vendorName"`
	Query    string `json:"query" validate:"required"`
}

type FilterOrdersSchema struct {
	Filters OrderFiltersSchema `json:"filters" validate:"required"`
}

type OrderFiltersSchema struct {
	BookNames  []string `json:"book_names"`
	Authors    []string `json:"authors"`
	Vendors    []string `json:"vendors"`
	NumOfItems []int    `json:"num_of_items"`
	TotalValue []int    `json:"total_value"`
	FromDate   string   `json:"from_date"`
	ToDate     string   `json:"to_date"`
}

type GetUsersOrdersSchema struct {
	UserID string `json:"user_id" validate:"required"`
}

type GetUserOrderSchema struct {
	UserID  string `json:"user_id" validate:"required"`
	OrderID string `json:"order_id" validate:"required"`
}
