package schemas

type AddReviewSchema struct {
	BookID  string `json:"bookID" validate:"required"`
	Comment string `json:"comment" validate:"required"`
	Rating  int    `json:"rating" validate:"required"`
}

type UpdateReviewSchema struct {
	ID      string `json:"id" validate:"required"`
	Comment string `json:"comment" validate:"required"`
	Rating  int    `json:"rating" validate:"required"`
}

type DeleteReviewSchema struct {
	ID string `json:"id" validate:"required"`
}

type DeleteReviewsSchema struct {
	IDs []string `json:"ids" validate:"required"`
}
