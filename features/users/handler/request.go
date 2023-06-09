package handler

import "mime/multipart"

type InputUpdateProfile struct {
	ID       uint                  `json:"id"`
	Name     string                `json:"name"`
	Email    string                `json:"email"`
	Phone    string                `json:"phone_number"`
	Pictures *multipart.FileHeader `json:"pictures"`
}
