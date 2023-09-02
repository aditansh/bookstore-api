package schemas

type RegisterUserSchema struct {
	Name     string `json:"name" validate:"required,min=3"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserSchema struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=3"`
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=20"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}
