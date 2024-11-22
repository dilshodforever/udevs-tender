package request_model

import "time"

type CreateUserReq struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type CreateBidReq struct {
	Price        float64 `json:"price"`
	DeliveryTime int     `json:"delivery_time"`
	Comments     string  `json:"comments"`
}

type CreateTenderReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Budget      float64   `json:"budget"`
}

type UpdateTenderReq struct {
	Status string `json:"status"`
}

type CreateNotificationReq struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}
