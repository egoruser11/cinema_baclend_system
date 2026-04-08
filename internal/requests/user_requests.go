package requests

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Age      *uint   `json:"age"`
}

type AddUserMoneyBalanceRequest struct {
}
