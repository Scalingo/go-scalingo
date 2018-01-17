package scalingo

// StaticTokenGenerator is an implementation of TokenGenerator which always return the same token.
// This token is provided in the GetStaticTokenGenerator function.
type StaticTokenGenerator struct {
	token  string
	client *Client
}

// GetStaticTokenGenerator returns a new StaticTokenGenerator. The only argument is the token that
// will always be returned by this generator.
func GetStaticTokenGenerator(token string) *StaticTokenGenerator {
	return &StaticTokenGenerator{
		token: token,
	}
}

// GetAccessToken always returns the configured token.
func (t *StaticTokenGenerator) GetAccessToken() (string, error) {
	return t.token, nil
}

// SetClient sets the client attribute of this token generator.
func (t *StaticTokenGenerator) SetClient(c *Client) {
	t.client = c
}
