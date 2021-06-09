package models

type (
	User struct {
		User         string `json:"user"`
		Role         string `json:"role"`
		Password     string `json:"password,omitempty"`
		HashPassword []byte `json:"hashpassword,omitempty"`
	}
)
