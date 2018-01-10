package scalingo

import (
	"crypto/tls"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client interface {
	AddonsList(app string) ([]*Addon, error)
	AddonProvision(app, addon, planID string) (AddonRes, error)
	AddonDestroy(app, addonID string) error
	AddonUpgrade(app, addonID, planID string) (AddonRes, error)

	AddonProvidersList() ([]*AddonProvider, error)
	AddonProviderPlansList(addon string) ([]*Plan, error)

	AppsList() ([]*App, error)
	AppsShow(appName string) (*App, error)
	AppsDestroy(name string, currentName string) error
	AppsRestart(app string, scope *AppsRestartParams) (*http.Response, error)
	AppsCreate(opts AppsCreateOpts) (*App, error)
	AppsStats(app string) (*AppStatsRes, error)
	AppsPs(app string) ([]ContainerType, error)
	AppsScale(app string, params *AppsScaleParams) (*http.Response, error)

	CollaboratorsList(app string) ([]Collaborator, error)
	CollaboratorAdd(app string, email string) (Collaborator, error)
	CollaboratorRemove(app string, id string) error

	DeploymentList(app string) ([]*Deployment, error)
	Deployment(app string, deploy string) (*Deployment, error)
	DeploymentLogs(deployURL string) (*http.Response, error)
	DeploymentStream(deployURL string) (*websocket.Conn, error)
	DeploymentsCreate(app string, params *DeploymentsCreateParams) (*Deployment, error)

	DomainsList(app string) ([]Domain, error)
	DomainsAdd(app string, d Domain) (Domain, error)
	DomainsRemove(app string, id string) error
	DomainsUpdate(app, id, cert, key string) (Domain, error)

	VariablesList(app string) (Variables, error)
	VariablesListWithoutAlias(app string) (Variables, error)
	VariableSet(app string, name string, value string) (*Variable, int, error)
	VariableMultipleSet(app string, variables Variables) (Variables, int, error)
	VariableUnset(app string, id string) error

	EventsList(app string, opts PaginationOpts) (Events, PaginationMeta, error)
	UserEventsList(opts PaginationOpts) (Events, PaginationMeta, error)

	KeysList() ([]Key, error)
	KeysAdd(name string, content string) (*Key, error)
	KeysDelete(id string) error

	Login(email, password string) (*LoginResponse, error)

	LogsArchivesByCursor(app string, cursor string) (*LogsArchivesResponse, error)
	LogsArchives(app string, page int) (*LogsArchivesResponse, error)
	LogsURL(app string) (*http.Response, error)
	Logs(logsURL string, n int, filter string) (*http.Response, error)

	NotificationPlatformsList() ([]*NotificationPlatform, error)
	NotificationPlatformByName(name string) ([]*NotificationPlatform, error)
	NotificationsList(app string) ([]*Notification, error)
	NotificationProvision(app, webHookURL string) (NotificationRes, error)
	NotificationUpdate(app, ID, webHookURL string) (NotificationRes, error)
	NotificationDestroy(app, ID string) error

	NotifiersList(app string) (Notifiers, error)
	NotifierProvision(app, notifierType string, params NotifierParams) (*Notifier, error)
	NotifierByID(app, ID string) (*Notifier, error)
	NotifierUpdate(app, ID, notifierType string, params NotifierParams) (*Notifier, error)
	NotifierDestroy(app, ID string) error

	OperationsShow(app string, opID string) (*Operation, error)

	SignUp(email, password string) error
	Self() (*User, error)
	UpdateUser(params UpdateUserParams) (*User, error)

	Run(opts RunOpts) (*RunRes, error)

	SourcesCreate() (*Source, error)

	TokensList() (Tokens, error)
	CreateToken(t Token) (Token, error)
	ShowToken(id int) (Token, error)
}

type clientImpl struct {
	TokenGenerator TokenGenerator
	Endpoint       string
	TLSConfig      *tls.Config
	APIVersion     string
	httpClient     HTTPClient
}

type ClientConfig struct {
	Timeout        time.Duration
	Endpoint       string
	TLSConfig      *tls.Config
	TokenGenerator TokenGenerator
}

func NewClient(cfg ClientConfig) Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = defaultEndpoint
	}
	if cfg.TLSConfig == nil {
		cfg.TLSConfig = &tls.Config{}
	}
	return &clientImpl{
		TokenGenerator: cfg.TokenGenerator,
		Endpoint:       cfg.Endpoint,
		APIVersion:     defaultAPIVersion,
		TLSConfig:      cfg.TLSConfig,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Proxy:           http.ProxyFromEnvironment,
				TLSClientConfig: cfg.TLSConfig,
			},
		},
	}
}

func (c *clientImpl) HTTPClient() HTTPClient {
	return c.httpClient
}
