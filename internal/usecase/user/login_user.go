package user

type LoginUserReq struct {
	Email    string `json:"email" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required"`
}
