package http

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	error
}

