package scalingo

type Client struct {
	APIToken string
}

func NewClient(token string) *Client {
	return &Client{APIToken: token}
}
