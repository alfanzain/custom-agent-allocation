package types

type AllocateAgentWebhookCandidateAgent struct {
	AvatarURL    *string `json:"avatar_url"`
	CreatedAt    string  `json:"created_at"`
	Email        string  `json:"email"`
	ForceOffline bool    `json:"force_offline"`
	ID           uint    `json:"id"`
	IsAvailable  bool    `json:"is_available"`
	IsVerified   bool    `json:"is_verified"`
	LastLogin    *string `json:"last_login"`
	Name         string  `json:"name"`
	SDKEmail     string  `json:"sdk_email"`
	SDKKey       string  `json:"sdk_key"`
	Type         int     `json:"type"`
	TypeAsString string  `json:"type_as_string"`
	UpdatedAt    string  `json:"updated_at"`
}

type AllocateAgentWebhookPayload struct {
	AppID          string                              `json:"app_id"`
	AvatarURL      string                              `json:"avatar_url"`
	CandidateAgent *AllocateAgentWebhookCandidateAgent `json:"candidate_agent"`
	Email          string                              `json:"email"`
	Extras         string                              `json:"extras"`
	IsNewSession   bool                                `json:"is_new_session"`
	IsResolved     bool                                `json:"is_resolved"`
	LatestService  *string                             `json:"latest_service"`
	Name           string                              `json:"name"`
	RoomID         string                              `json:"room_id"`
	Source         string                              `json:"source"`
}
