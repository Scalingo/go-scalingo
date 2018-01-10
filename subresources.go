package scalingo

import "gopkg.in/errgo.v1"

func (c *clientImpl) subresourceGet(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *clientImpl) subresourceList(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/" + subresource,
		Params:   payload,
	}, data)
}

func (c *clientImpl) subresourceAdd(app, subresource string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/" + subresource,
		Expected: Statuses{201},
		Params:   payload,
	}, data)
}

func (c *clientImpl) subresourceDelete(app string, subresource string, id string) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Expected: Statuses{204},
	}, nil)
}

func (c *clientImpl) subresourceUpdate(app, subresource, id string, payload, data interface{}) error {
	return c.doSubresourceRequest(&APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *clientImpl) doSubresourceRequest(req *APIRequest, data interface{}) error {
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
