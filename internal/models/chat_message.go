package models

import "time"

type ChatMessage struct {
	ID             int       `json:"id"`
	Channel        string    `json:"channel"`
	AuthorID       int       `json:"author_id"`
	AuthorName     string    `json:"author_name"`
	AuthorRole     string    `json:"author_role"`
	Content        string    `json:"content"`
	AttachmentsURL *string   `json:"attachments_url,omitempty"`
	SentAt         time.Time `json:"sent_at"`
}
