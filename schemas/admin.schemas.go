package schemas

type PromoteDeactivateDeleteAccountsSchema struct {
	Usernames []string `json:"usernames" validate:"required"`
}

type GetUserCartOrdersSchema struct {
	UserID string `json:"user_id" validate:"required"`
}

type ModifyUserCartSchema struct {
	UserID   string `json:"user_id" validate:"required"`
	BookID   string `json:"book_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type ClearUserCartSchema struct {
	UserID string `json:"user_id" validate:"required"`
}

type ApproveVendorSchema struct {
	ID string `json:"id" validate:"required"`
}

type FlagUserVendorSchema struct {
	Username string `json:"username" validate:"required"`
}
