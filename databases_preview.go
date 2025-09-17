package scalingo

import (
	"context"
	stderrors "errors"

	"github.com/Scalingo/go-utils/errors/v2"
)

const databasesResource = "databases"

var ErrDatabaseNotFound = stderrors.New("database not found")

type PreviewAPI interface {
	DatabasesPreviewService
}

type DatabasesPreviewService interface {
	DatabaseCreate(ctx context.Context, params DatabaseCreateParams) (DedicatedDatabase, error)
	DatabasesList(ctx context.Context) ([]DedicatedDatabase, error)
	DatabaseShow(ctx context.Context, appID string) (DedicatedDatabaseShowResponse, error)
	DatabaseDestroy(ctx context.Context, appID string) error
}

var _ DatabasesPreviewService = (*PreviewClient)(nil)

type DatabaseCreateParams struct {
	AddonProviderID string `json:"addon_provider_id"`
	PlanID          string `json:"plan_id"`
	DatabaseName    string `json:"database_name"`
	ProjectID       string `json:"project_id,omitempty"`
}

type DedicatedDatabase struct {
	App   App   `json:"app"`
	Addon Addon `json:"addon"`
}

type DedicatedDatabaseShowResponse struct {
	App      App      `json:"app"`
	Addon    Addon    `json:"addon"`
	Database Database `json:"database"`
}

type PreviewClient struct {
	parent *Client
}

func NewPreviewClient(parent *Client) PreviewClient {
	return PreviewClient{
		parent: parent,
	}
}

func (c *PreviewClient) DatabaseCreate(ctx context.Context, params DatabaseCreateParams) (DedicatedDatabase, error) {
	var res DedicatedDatabase
	err := c.parent.ScalingoAPI().ResourceAdd(ctx, databasesResource, params, &res)
	if err != nil {
		return res, errors.Wrap(ctx, err, "create dedicated database")
	}
	return res, nil
}

func (c *PreviewClient) DatabasesList(ctx context.Context) ([]DedicatedDatabase, error) {
	var res []DedicatedDatabase

	err := c.parent.ScalingoAPI().ResourceList(ctx, databasesResource, nil, &res)
	if err != nil {
		return res, errors.Wrap(ctx, err, "list dedicated databases")
	}
	return res, nil
}

// DatabaseShow currently uses appID as the dedicated database identifier.
func (c *PreviewClient) DatabaseShow(ctx context.Context, appID string) (DedicatedDatabaseShowResponse, error) {
	var res DedicatedDatabaseShowResponse

	dedicatedDatabase, err := c.searchDatabase(ctx, appID)
	if err != nil {
		return res, errors.Wrap(ctx, err, "search dedicated database")
	}

	database, err := c.parent.DatabaseShow(ctx, dedicatedDatabase.App.ID, dedicatedDatabase.Addon.ID)
	if err != nil {
		return res, errors.Wrap(ctx, err, "show dedicated database")
	}

	res.App = dedicatedDatabase.App
	res.Addon = dedicatedDatabase.Addon
	res.Database = database

	return res, nil
}

// DatabaseDestroy currently uses appID as the dedicated database identifier.
func (c *PreviewClient) DatabaseDestroy(ctx context.Context, appID string) error {
	dedicatedDatabase, err := c.searchDatabase(ctx, appID)
	if err != nil {
		return errors.Wrap(ctx, err, "search dedicated database")
	}

	appName := dedicatedDatabase.App.Name
	err = c.parent.AppsDestroy(ctx, appName, appName)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy dedicated database app")
	}
	return nil
}

// searchDatabase performs a linear search through DatabasesList method result.
func (c *PreviewClient) searchDatabase(ctx context.Context, appID string) (DedicatedDatabase, error) {
	var res DedicatedDatabase

	databases, err := c.DatabasesList(ctx)
	if err != nil {
		return res, errors.Wrap(ctx, err, "search dedicated database")
	}

	for _, dedicatedDatabase := range databases {
		if dedicatedDatabase.App.ID == appID {
			return dedicatedDatabase, nil
		}
	}
	return res, ErrDatabaseNotFound
}
