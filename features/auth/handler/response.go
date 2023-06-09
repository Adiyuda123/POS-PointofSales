package handler

type LoginResponse struct {
	ID    uint   `json:"user_id"`
	Name  string `json:"user_name"`
	Token string `json:"token"`
}

type RegisterResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone_number"`
	Pictures string `json:"pictures"`
}
