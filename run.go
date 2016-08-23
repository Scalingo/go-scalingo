package scalingo

import (
	"net/http"
	"strings"
)

type RunOpts struct {
	App      string
	Command  []string
	Env      map[string]string
	Size     string
	Detached bool
}

func (c *Client) Run(opts RunOpts) (*http.Response, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "POST",
		Endpoint: "/apps/" + opts.App + "/run",
		Params: map[string]interface{}{
			"command":  strings.Join(opts.Command, " "),
			"env":      opts.Env,
			"size":     opts.Size,
			"detached": opts.Detached,
		},
	}
	return req.Do()
}
