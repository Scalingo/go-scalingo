package scalingo

import "gopkg.in/errgo.v1"

type Collaborator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type CollaboratorsRes struct {
	Collaborators []Collaborator `json:"collaborators"`
}

type CollaboratorRes struct {
	Collaborator Collaborator `json:"collaborator"`
}

func (c *clientImpl) CollaboratorsList(app string) ([]Collaborator, error) {
	var collaboratorsRes CollaboratorsRes
	err := c.subresourceList(app, "collaborators", nil, &collaboratorsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return collaboratorsRes.Collaborators, nil
}

func (c *clientImpl) CollaboratorAdd(app string, email string) (Collaborator, error) {
	var collaboratorRes CollaboratorRes
	err := c.subresourceAdd(app, "collaborators", CollaboratorRes{
		Collaborator: Collaborator{Email: email},
	}, &collaboratorRes)
	if err != nil {
		return Collaborator{}, errgo.Mask(err)
	}
	return collaboratorRes.Collaborator, nil
}

func (c *clientImpl) CollaboratorRemove(app string, id string) error {
	return c.subresourceDelete(app, "collaborators", id)
}
