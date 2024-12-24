package types

type MarkAsResolvedWebhookCustomer struct {
	AdditionalInfo []interface{} `json:"additional_info"`
	Avatar         string        `json:"avatar"`
	Name           string        `json:"name"`
	UserID         string        `json:"user_id"`
}

type MarkAsResolvedWebhookResolvedBy struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsAvailable bool   `json:"is_available"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type MarkAsResolvedWebhookService struct {
	FirstCommentID string  `json:"first_comment_id"`
	ID             int     `json:"id"`
	IsResolved     bool    `json:"is_resolved"`
	LastCommentID  string  `json:"last_comment_id"`
	Notes          *string `json:"notes"`
	RoomID         string  `json:"room_id"`
	Source         string  `json:"source"`
}

type MarkAsResolvedWebhookPayload struct {
	Customer   MarkAsResolvedWebhookCustomer   `json:"customer"`
	ResolvedBy MarkAsResolvedWebhookResolvedBy `json:"resolved_by"`
	Service    MarkAsResolvedWebhookService    `json:"service"`
}
