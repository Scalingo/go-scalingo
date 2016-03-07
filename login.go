package scalingo

import (
	"net/http"
)

type LoginError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	AuthenticationToken string `json:"authentication_token"`
	User                *User  `json:"user"`
}

func (err *LoginError) Error() string {
	return err.Message
}

func Login(email, password string) (*http.Response, error) {
	req := &APIRequest{
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/users/sign_in",
		Expected: Statuses{201, 401},
		Params: map[string]interface{}{
			"user": map[string]string{
				"login":    email,
				"password": password,
			},
		},
	}
	return req.Do()
}
