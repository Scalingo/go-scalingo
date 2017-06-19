package scalingo

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/Scalingo/cli/debug"
)

type Notifier struct {
	ID             string                 `json:"id"`
	Active         bool                   `json:"active"`
	Name           string                 `json:"name"`
	Type           NotifierType           `json:"type"`
	SendAllEvents  bool                   `json:"send_all_events"`
	SelectedEvents []string               `json:"selected_events"`
	TypeData       map[string]interface{} `json:"-"`
	RawTypeData    json.RawMessage        `json:"type_data"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type NotifierType string

const (
	NotifierWebhook NotifierType = "webhook"
	NotifierSlack                = "slack"
)

type DetailedNotifier interface {
	// fmt.Stringer
	GetNotifier() *Notifier
	GetID() string
	GetName() string
	GetType() NotifierType
	GetSendAllEvents() bool
	GetSelectedEvents() []string
	IsActive() bool
	// PrintableType() string
	When() string
	// Who() string
	TypeDataPtr() interface{}
	TypeDataString() string
}

type Notifiers []DetailedNotifier

// DetailedNotifier implementation
func (not *Notifier) GetNotifier() *Notifier {
	return not
}

func (not *Notifier) GetID() string {
	return not.ID
}

func (not *Notifier) GetName() string {
	return not.Name
}

func (not *Notifier) GetType() NotifierType {
	return not.Type
}

func (not *Notifier) GetSendAllEvents() bool {
	return not.SendAllEvents
}

func (not *Notifier) GetSelectedEvents() []string {
	return not.SelectedEvents
}

func (not *Notifier) IsActive() bool {
	return not.Active
}

func (not *Notifier) When() string {
	return not.UpdatedAt.Format("Mon Jan 02 2006 15:04:05")
}

func (not *Notifier) TypeDataPtr() interface{} {
	return &not.TypeData
}

func (not *Notifier) TypeDataString() string {
	return "unknow notifier type"
}

// Webhook
type NotifierWebhookType struct {
	Notifier
	TypeData NotifierWebhookTypeData `json:"type_data"`
}

type NotifierWebhookTypeData struct {
	WebhookURL string `json:"webhook_url"`
}

func (e *NotifierWebhookType) TypeDataPtr() interface{} {
	return &e.TypeData
}

func (not *NotifierWebhookType) TypeDataString() string {
	return fmt.Sprintf("- webhook url: %s", not.TypeData.WebhookURL)
}

// Slack
type NotifierSlackType struct {
	Notifier
	TypeData NotifierSlackTypeData `json:"type_data"`
}

type NotifierSlackTypeData struct {
	WebhookURL string `json:"webhook_url"`
}

func (e *NotifierSlackType) TypeDataPtr() interface{} {
	return &e.TypeData
}

func (not *NotifierSlackType) TypeDataString() string {
	return fmt.Sprintf("- webhook url: %s", not.TypeData.WebhookURL)
}

func (pnot *Notifier) Specialize() DetailedNotifier {
	var detailedNotifier DetailedNotifier
	notifier := *pnot
	switch notifier.Type {
	case NotifierWebhook:
		detailedNotifier = &NotifierWebhookType{Notifier: notifier}
	case NotifierSlack:
		detailedNotifier = &NotifierSlackType{Notifier: notifier}
	default:
		return pnot
	}
	err := json.Unmarshal(pnot.RawTypeData, detailedNotifier.TypeDataPtr())
	if err != nil {
		debug.Printf("error reading the data: %+v\n", err)
		return pnot
	}
	return detailedNotifier
}

func NewNotifier(notifierType string, m map[string]interface{}) (notifier interface{}) {
	switch notifierType {
	case "webhook":
		notifier = &NotifierWebhookType{}
		fillStruct(&NotifierWebhookType{}, m)
	case "slack":
		notifier = &NotifierSlackType{}
		fillStruct(&NotifierSlackType{}, m)
	}
	return
}

// private
func fillStruct(str interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(str, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}
