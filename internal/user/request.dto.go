package user

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30`
	Password string `json:"password" binding:"required,min=6,max=100`
}
