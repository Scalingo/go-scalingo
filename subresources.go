package scalingo

import (
	"crypto/tls"

	"gopkg.in/errgo.v1"
)

type backend interface {
	subresourceGet(app, subresource, id string, payload, data interface{}) error
	subresourceList(app, subresource string, payload, data interface{}) error
	subresourceAdd(app, subresource string, payload, data interface{}) error
	subresourceDelete(app string, subresource string, id string) error
	subresourceUpdate(app, subresource, id string, payload, data interface{}) error
	doSubresourceRequest(req *APIRequest, data interface{}) error
}

type backendConfiguration struct {
	TokenGenerator TokenGenerator
	Endpoint       string
	TLSConfig      *tls.Config
	APIVersion     string
	httpClient     HTTPClient
}

func (c *backendConfiguration) subresourceGet(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *backendConfiguration) subresourceList(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource,
		Params:   payload,
	}, data)
}

func (c *backendConfiguration) subresourceAdd(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/" + subresource,
		Expected: Statuses{201},
		Params:   payload,
	}, data)
}

func (c *backendConfiguration) subresourceDelete(app string, subresource string, id string) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Expected: Statuses{204},
	}, nil)
}

func (c *backendConfiguration) subresourceUpdate(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *backendConfiguration) doSubresourceRequest(req *APIRequest, data interface{}) error {
	req.Client = c
	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	if data == nil {
		return nil
	}

	err = ParseJSON(res, data)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
