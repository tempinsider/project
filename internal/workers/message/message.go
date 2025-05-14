package message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"insider-mert/internal/models"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	RetrieveItemsFromDatabase = `
		SELECT id, message_content, phone_number
		FROM messages
		WHERE status = $1
		LIMIT 2
		FOR UPDATE SKIP LOCKED
	`
	UpdateMessageStatus = `
        UPDATE messages 
        SET status = $1, external_message_id = $2
        WHERE id = $3
    `
	MaxWorkers = 10
)

var (
	WebhookUrl string
)

func init() {
	WebhookUrl = os.Getenv("WEBHOOK_SITE_URL")
}

type MessageResult struct {
	Message *models.Message
	Error   error
}

func SendMessage(conn *pgx.Conn) {
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Default().Printf("error during sending message, err : %v", err.Error())
		return
	}
	// we cant use defer tx.Rollback(ctx) here, because if we commit data, it could literally try rollback for no reason.

	var items []*models.Message
	rows, err := tx.Query(ctx, RetrieveItemsFromDatabase, models.MessageStatusOpen)
	if err != nil {
		tx.Rollback(ctx)
		log.Default().Printf("error receiving items from db, err : %v", err.Error())
		return
	}

	for rows.Next() {
		var message models.Message
		err = rows.Scan(
			&message.ID,
			&message.MessageContent,
			&message.PhoneNumber,
		)
		if err != nil {
			tx.Rollback(ctx)
			log.Default().Printf("error scanning row, err: %v", err.Error())
			return
		}
		items = append(items, &message)
	}

	if len(items) == 0 {
		tx.Rollback(ctx)
		return
	}

	jobs := make(chan *models.Message, len(items))
	successCh := make(chan *models.Message, len(items))

	// Since we have 10 workers in config but actually have 2 items, we need to reduce workers. but, it may still helpfull if we have 2 workers but only 1 record in db.
	workerCount := MaxWorkers
	if len(items) < MaxWorkers {
		workerCount = len(items)
	}

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, successCh, &wg)
	}

	for _, msg := range items {
		jobs <- msg
	}
	close(jobs)

	wg.Wait()
	close(successCh)

	successCount := 0
	for msg := range successCh {
		_, err := tx.Exec(ctx, UpdateMessageStatus,
			models.MessageStatusSent, msg.ExternalMessageID, msg.ID)
		if err != nil {
			log.Default().Printf("error updating message status for ID %v: %v",
				msg.ID, err.Error())
			continue
		}
		successCount++
	}

	if successCount > 0 {
		err = tx.Commit(ctx)
		if err != nil {
			log.Default().Printf("error committing transaction: %v", err.Error())
			return
		}
	} else {
		log.Default().Printf("no messages were sent successfully, rolling back")
		tx.Rollback(ctx)
	}
}

func worker(jobs <-chan *models.Message,
	successCh chan<- *models.Message, wg *sync.WaitGroup) {

	defer wg.Done()

	for msg := range jobs {
		messageID, err := sendPostRequest(msg)

		if err != nil {
			continue
		}

		if messageID != "" {
			msg.ExternalMessageID = &messageID
			successCh <- msg
		}
	}
}

func sendPostRequest(message *models.Message) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	payload, err := json.Marshal(map[string]string{
		"to":      message.PhoneNumber,
		"content": message.MessageContent,
	})
	if err != nil {
		return "", err
	}

	resp, err := client.Post(
		WebhookUrl,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 202 {
		rawData, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", nil // not sure what to tbh, server accepted it but we somehow cant parse response data?
		}

		type responseData struct {
			MessageID string `json:"messageId"`
		}

		var parsedData responseData
		err = json.Unmarshal(rawData, &parsedData)
		if err != nil {
			return "", nil // not sure what to tbh, server accepted it but we somehow cant parse response data?
		}

		return parsedData.MessageID, nil
	}
	return "", fmt.Errorf("service returned status code: %d", resp.StatusCode)
}
