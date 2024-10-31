package devhub

type Database struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Adapter        string `json:"adapter"`
	Hostname       string `json:"hostname"`
	Database       string `json:"database"`
	Ssl            bool   `json:"ssl"`
	Cacertfile     string `json:"cacertfile"`
	Keyfile        string `json:"keyfile"`
	Certfile       string `json:"certfile"`
	RestrictAccess bool   `json:"restrict_access"`
}
