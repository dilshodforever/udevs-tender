package response_model

type ProfileRes struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginRes struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}
