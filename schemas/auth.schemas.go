package schemas

type LoginUserSchema struct {
	Username string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginVendorSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponseSchema struct {
	RefreshToken string `json:"refreshToken"`
}

type VerifyOTPSchema struct {
	OTP string `json:"otp" validate:"required,min=6,max=6"`
}

type ResendOTPSchema struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordSchema struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordSchema struct {
	Password string `json:"password" validate:"required,min=6"`
}

type UpdatePasswordSchema struct {
	OldPassword string `json:"oldPassword" validate:"required,min=6"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

type RefreshTokenSchema struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponseSchema struct {
	AccessToken string `json:"accessToken"`
}

type LogoutSchema struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
