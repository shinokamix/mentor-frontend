package models

type ReviewEvent struct {
	Action string  `json:"action"` // created/updated/deleted
	ID     int64   `json:"id"`
	Email  string  `json:"email"`
	Score  float32 `json:"score"`
}
