package scalingo

type StaticTokenGenerator struct {
	token  string
	client *Client
}

func (c *Client) GetStaticTokenGenerator(token string) *StaticTokenGenerator {
	return &StaticTokenGenerator{
		token: token,
	}
}

func (t *StaticTokenGenerator) GetAccessToken() (string, error) {
	return t.token, nil
}

func (t *StaticTokenGenerator) SetClient(c *Client) {
	t.client = c
}
