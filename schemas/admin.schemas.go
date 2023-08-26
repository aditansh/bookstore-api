package schemas

type DeactivateDeleteAccountsSchema struct {
	IDs []string `json:"ids" validate:"required"`
}

type ModifyUserCartSchema struct {
	UserID   string `json:"user_id" validate:"required"`
	BookID   string `json:"book_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}
