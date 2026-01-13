package model

// RequestCreateUser represents the request body for user registration.
type RequestCreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ResponseLoginUser represents the response body after successful login or registration.
type ResponseLoginUser struct {
	AccessToken string `json:"access_token,omitempty"`
	ID          string `json:"id"`
	Username    string `json:"username"`
}
