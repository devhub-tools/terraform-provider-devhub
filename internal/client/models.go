package devhub

type Database struct {
	Id                   string               `json:"id"`
	Name                 string               `json:"name"`
	Adapter              string               `json:"adapter"`
	Hostname             string               `json:"hostname"`
	Database             string               `json:"database"`
	Ssl                  bool                 `json:"ssl"`
	Cacertfile           string               `json:"cacertfile"`
	Keyfile              string               `json:"keyfile"`
	Certfile             string               `json:"certfile"`
	RestrictAccess       bool                 `json:"restrict_access"`
	EnableDataProtection bool                 `json:"enable_data_protection"`
	Group                string               `json:"group"`
	SlackWebhookURL      string               `json:"slack_webhook_url"`
	SlackChannel         string               `json:"slack_channel"`
	AgentId              string               `json:"agent_id"`
	Credentials          []DatabaseCredential `json:"credentials"`
}

type DatabaseCredential struct {
	Id                string `json:"id"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ReviewsRequired   int    `json:"reviews_required"`
	DefaultCredential bool   `json:"default_credential"`
}
