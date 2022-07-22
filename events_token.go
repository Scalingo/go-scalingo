package scalingo

type EventNewTokenTypeData struct {
	TokenName string `json:"token_name"`
	TokenId   string `json:"token_id"`
}

type EventNewTokenType struct {
	Event
	TypeData EventNewTokenTypeData `json:"type_data"`
}

func (ev *EventNewTokenType) String() string {
	return "New token created"
}

type EventRegenerateTokenTypeData struct {
	TokenName string `json:"token_name"`
	TokenId   string `json:"token_id"`
}

type EventRegenerateTokenType struct {
	Event
	TypeData EventRegenerateTokenTypeData `json:"type_data"`
}

func (ev *EventRegenerateTokenType) String() string {
	return "Token regenerated"
}

type EventDeleteTokenTypeData struct {
	TokenName string `json:"token_name"`
	TokenId   string `json:"token_id"`
}

type EventDeleteTokenType struct {
	Event
	TypeData EventDeleteTokenTypeData `json:"type_data"`
}

func (ev *EventDeleteTokenType) String() string {
	return "Token deleted"
}
