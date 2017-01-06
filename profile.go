package scalingo

import (
	"encoding/json"

	"gopkg.in/errgo.v1"
)

type UpdateUserParams struct {
	StopFreeTrial string `json:"stop_free_trial,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
}

type UpdateUserResponse struct {
	User *User `json:"user"`
}

func (c *Client) UpdateUser(params UpdateUserParams) (*User, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "PATCH",
		Endpoint: "/account/profile",
		Params: map[string]interface{}{
			"user": map[string]interface{}{
				"stop_free_trial": true,
			},
		},
		Expected: Statuses{200},
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var u *UpdateUserResponse
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return u.User, nil
}
