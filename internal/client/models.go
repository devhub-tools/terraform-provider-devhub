package devhub

type Database struct {
	Id             string               `json:"id"`
	Name           string               `json:"name"`
	Adapter        string               `json:"adapter"`
	Hostname       string               `json:"hostname"`
	Database       string               `json:"database"`
	Ssl            bool                 `json:"ssl"`
	Cacertfile     string               `json:"cacertfile"`
	Keyfile        string               `json:"keyfile"`
	Certfile       string               `json:"certfile"`
	RestrictAccess bool                 `json:"restrict_access"`
	Group          string               `json:"group"`
	SlackChannel   string               `json:"slack_channel"`
	AgentId        string               `json:"agent_id"`
	Credentials    []DatabaseCredential `json:"credentials"`
}

type DatabaseCredential struct {
	Id                string `json:"id"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ReviewsRequired   int    `json:"reviews_required"`
	DefaultCredential bool   `json:"default_credential"`
}

type TerradeskWorkspace struct {
	Id                    string            `json:"id"`
	Name                  string            `json:"name"`
	Repository            string            `json:"repository"`
	InitArgs              string            `json:"init_args"`
	Path                  string            `json:"path"`
	RunPlansAutomatically bool              `json:"run_plans_automatically"`
	RequiredApprovals     int               `json:"required_approvals"`
	DockerImage           string            `json:"docker_image"`
	CpuRequests           string            `json:"cpu_requests"`
	MemoryRequests        string            `json:"memory_requests"`
	AgentId               string            `json:"agent_id"`
	WorkloadIdentity      *WorkloadIdentity `json:"workload_identity"`
	EnvVars               []EnvVar          `json:"env_vars"`
	Secrets               []Secret          `json:"secrets"`
}

type WorkloadIdentity struct {
	Enabled             bool   `json:"enabled"`
	ServiceAccountEmail string `json:"service_account_email"`
	Provider            string `json:"provider"`
}

type EnvVar struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Secret struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
