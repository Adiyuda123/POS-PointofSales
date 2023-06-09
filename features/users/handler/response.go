package handler

type GetUserByIdResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Pictures string `json:"pictures"`
}
