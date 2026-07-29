package devhub

type Database struct {
	Id                   string               `json:"id"`
	Name                 string               `json:"name"`
	Adapter              string               `json:"adapter"`
	Hostname             string               `json:"hostname"`
	Port                 *int64               `json:"port"`
	Database             string               `json:"database"`
	Ssl                  bool                 `json:"ssl"`
	Cacertfile           string               `json:"cacertfile"`
	Keyfile              string               `json:"keyfile"`
	Certfile             string               `json:"certfile"`
	VerifyServerHostname bool                 `json:"verify_server_hostname"`
	RestrictAccess       bool                 `json:"restrict_access"`
	Group                string               `json:"group"`
	SlackChannel         string               `json:"slack_channel"`
	AgentId              string               `json:"agent_id"`
	AiEnabled            bool                 `json:"ai_enabled"`
	AiMaxRows            int64                `json:"ai_max_rows"`
	Credentials          []DatabaseCredential `json:"credentials"`
}

type DatabaseCredential struct {
	Id                string `json:"id"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ReviewsRequired   int    `json:"reviews_required"`
	DefaultCredential bool   `json:"default_credential"`
	AiAllowed         bool   `json:"ai_allowed"`
	Timeout           *int64 `json:"timeout"`
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

type Workflow struct {
	Id                 string             `json:"id"`
	Name               string             `json:"name"`
	CronSchedule       string             `json:"cron_schedule"`
	Group              string             `json:"group"`
	TriggerLinearLabel TriggerLinearLabel `json:"trigger_linear_label"`
	Inputs             []WorkflowInput    `json:"inputs"`
	Steps              []WorkflowStep     `json:"steps"`
}

type TriggerLinearLabel struct {
	Name string `json:"name"`
}

type WorkflowInput struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Type        string `json:"type"` // string, float, integer, boolean
	Required    bool   `json:"required"`
}

type WorkflowStep struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
	// Action fields
	Action      *WorkflowStepAction `json:"action"`
	Permissions []*Permission       `json:"permissions"`
}

type WorkflowStepAction struct {
	Type string `json:"__type__"`
	// ApiAction fields
	Endpoint           string                        `json:"endpoint"`
	Method             string                        `json:"method"`
	Headers            []WorkflowStepActionApiHeader `json:"headers"`
	Body               string                        `json:"body"`
	ExpectedStatusCode int64                         `json:"expected_status_code"`
	IncludeDevhubJwt   bool                          `json:"include_devhub_jwt"`
	// ApprovalAction fields
	ReviewsRequired int `json:"reviews_required"`
	// QueryAction fields
	Timeout      int    `json:"timeout"`
	Query        string `json:"query"`
	CredentialId string `json:"credential_id"`
	// SlackAction fields
	SlackChannel string `json:"slack_channel"`
	Message      string `json:"message"`
	LinkText     string `json:"link_text"`
	// SlackReplyAction fields
	ReplyToStepName string `json:"reply_to_step_name"`
	// ConditionAction fields
	Condition string `json:"condition"`
	WhenFalse string `json:"when_false"`
}

type WorkflowStepActionApiHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Dashboard struct {
	Id               string           `json:"id"`
	Name             string           `json:"name"`
	RestrictedAccess bool             `json:"restricted_access"`
	Panels           []DashboardPanel `json:"panels"`
}

type DashboardPanel struct {
	Id      string                 `json:"id"`
	Title   string                 `json:"title"`
	Inputs  []DashboardPanelInput  `json:"inputs"`
	Actions []DashboardPanelAction `json:"actions"`
	Details *DashboardPanelDetails `json:"details"`
}

type DashboardPanelInput struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

type DashboardPanelAction struct {
	Id            string                             `json:"id"`
	Label         string                             `json:"label"`
	WorkflowId    string                             `json:"workflow_id"`
	InputMappings []DashboardPanelActionInputMapping `json:"input_mappings"`
}

type DashboardPanelActionInputMapping struct {
	WorkflowInputKey string `json:"workflow_input_key"`
	Column           string `json:"column"`
	Prompt           string `json:"prompt"`
	PromptUser       bool   `json:"prompt_user"`
}

type DashboardPanelDetails struct {
	Type string `json:"__type__"`
	// QueryPanel fields
	Query        string `json:"query"`
	CredentialId string `json:"credential_id"`
}

type Permission struct {
	Id                 string `json:"id"`
	Permission         string `json:"permission"`
	RoleId             string `json:"role_id"`
	OrganizationUserId string `json:"organization_user_id"`
}

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Managed     bool   `json:"managed"`
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
