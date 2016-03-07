package scalingo

type Client struct {
	APIToken   string
	Endpoint   string
	APIVersion string
}

func NewClient(token string) *Client {
	return &Client{APIToken: token, Endpoint: defaultEndpoint, APIVersion: defaultAPIVersion}
}
