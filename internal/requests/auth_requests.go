package requests

type LoginRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceInfo string `json:"device_info"`
}
type RegisterRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Age        uint   `json:"age"`
	DeviceInfo string `json:"device_info"`
}
