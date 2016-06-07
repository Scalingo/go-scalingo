package scalingo

type Client struct {
	APIToken   string
	Endpoint   string
	APIVersion string
}

type ClientConfig struct {
	Endpoint string
	APIToken string
}

func NewClient(cfg ClientConfig) *Client {
	if cfg.Endpoint == "" {
		cfg.Endpoint = defaultEndpoint
	}
	return &Client{APIToken: cfg.APIToken, Endpoint: cfg.Endpoint, APIVersion: defaultAPIVersion}
}
