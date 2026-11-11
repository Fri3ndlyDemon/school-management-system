package repository

import (
	"database/sql"
	"school-management-system/internal/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// fetch chat messages for a specific channel in chronological order
func (r *StudentRepository) GetChannelMessages(channel string) ([]models.ChatMessage, error) {
	query := `
		SELECT cm.id, cm.channel, cm.author_id, u.name, u.role, cm.content, cm.attachments_url, cm.sent_at
		FROM chat_messages cm
		JOIN users u ON cm.author_id = u.id
		WHERE cm.channel = $1
		ORDER BY cm.sent_at ASC
	`
	rows, err := r.db.Query(query, channel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		var attachmentsURL sql.NullString
		if err := rows.Scan(&m.ID, &m.Channel, &m.AuthorID, &m.AuthorName, &m.AuthorRole, &m.Content, &attachmentsURL, &m.SentAt); err != nil {
			return nil, err
		}
		if attachmentsURL.Valid {
			m.AttachmentsURL = &attachmentsURL.String
		}
		msgs = append(msgs, m)
	}

	// check for any iteration errors encountered during rows scanning
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return msgs, nil
}

// insert a new text message into a designated chat channel
func (r *StudentRepository) SaveMessage(channel string, authorID int, content string) error {
	query := `INSERT INTO chat_messages (channel, author_id, content) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, channel, authorID, content)
	return err
}
