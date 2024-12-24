package responses

type QiscusAllocateAgentResponse struct {
	Data struct {
		Agent struct {
			ID           int     `json:"id"`
			Name         string  `json:"name"`
			SdkEmail     string  `json:"sdk_email"`
			Email        string  `json:"email"`
			SdkKey       string  `json:"sdk_key"`
			Type         int     `json:"type"`
			IsAvailable  bool    `json:"is_available"`
			AvatarURL    *string `json:"avatar_url"`
			IsVerified   bool    `json:"is_verified"`
			ForceOffline bool    `json:"force_offline"`
			Count        int     `json:"count"`
		} `json:"agent"`
	} `json:"data"`
}

type QiscusAssignAgentResponse struct {
	Data struct {
		AddedAgent struct {
			ID                  int      `json:"id"`
			Name                string   `json:"name"`
			Email               string   `json:"email"`
			AuthenticationToken string   `json:"authentication_token"`
			CreatedAt           string   `json:"created_at"`
			UpdatedAt           string   `json:"updated_at"`
			SDKEmail            string   `json:"sdk_email"`
			SDKKey              string   `json:"sdk_key"`
			IsAvailable         bool     `json:"is_available"`
			Type                int      `json:"type"`
			AvatarURL           string   `json:"avatar_url"`
			AppID               int      `json:"app_id"`
			IsVerified          bool     `json:"is_verified"`
			NotificationsRoomID string   `json:"notifications_room_id"`
			BubbleColor         string   `json:"bubble_color"`
			QismoKey            *string  `json:"qismo_key"`
			TypeAsString        string   `json:"type_as_string"`
			AssignedRules       []string `json:"assigned_rules"`
		} `json:"added_agent"`
	} `json:"data"`
}
