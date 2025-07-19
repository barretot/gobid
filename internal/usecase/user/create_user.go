package user

type CreateUserReq struct {
	UserName string `json:"user_name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Bio      string `json:"bio" validate:"required"`
}
