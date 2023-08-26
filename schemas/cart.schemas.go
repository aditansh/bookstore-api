package schemas

// type AddToCartSchema struct {
// 	BookID   string `json:"book_id" validate:"required"`
// 	Quantity int    `json:"quantity" validate:"required"`
// }

type RemoveFromCart struct {
	BookID string `json:"book_id" validate:"required"`
}

type ModifyCartSchema struct {
	BookID   string `json:"book_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type CheckoutSchema struct {
	Address string `json:"address" validate:"required"`
}
