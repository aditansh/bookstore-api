package schemas

type PromoteDeactivateDeleteAccountsSchema struct {
	Usernames []string `json:"usernames" validate:"required"`
}

type DeactivateDeleteVendorsSchema struct {
	Emails []string `json:"emails" validate:"required"`
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
	Email string `json:"email" validate:"email,required"`
}

type FlagUserVendorSchema struct {
	Email string `json:"email" validate:"email,required"`
}
