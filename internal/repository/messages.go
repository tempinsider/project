package repository

import (
	"context"
	"insider-mert/internal/models"
	"log"

	"github.com/jackc/pgx/v5"
)

const (
	PerPage               = 10
	ListSentMessagesQuery = `
		SELECT id,updated_at,external_message_id
		FROM messages
		WHERE status = $1
		ORDER BY updated_at
		LIMIT $2
		OFFSET $3
	`
)

type MessagesRepository struct {
	conn *pgx.Conn
}

func NewMessagesRepository(conn *pgx.Conn) *MessagesRepository {
	return &MessagesRepository{conn: conn}
}

func (m *MessagesRepository) ListSentMessages(ctx context.Context, page int) ([]*models.Message, error) {
	tx, err := m.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	offset := PerPage * page
	if page == 1 {
		offset = 0
	}

	var items []*models.Message
	rows, err := tx.Query(ctx, ListSentMessagesQuery, models.MessageStatusSent, PerPage, offset)
	if err != nil {
		tx.Rollback(ctx)
		log.Default().Printf("error receiving items from db, err : %v", err.Error())
		return nil, err
	}

	for rows.Next() {
		var message models.Message
		err = rows.Scan(
			&message.ID,
			&message.UpdatedAt,
			&message.ExternalMessageID,
		)
		if err != nil {
			tx.Rollback(ctx)
			log.Default().Printf("error scanning row, err: %v", err.Error())
			return nil, err
		}
		items = append(items, &message)
	}

	return items, nil
}
