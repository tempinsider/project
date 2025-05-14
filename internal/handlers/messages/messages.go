package messages

import (
	"encoding/json"
	"fmt"
	"insider-mert/internal/repository"
	"net/http"
	"strconv"
)

type MessagesHandler struct {
	messagesRepository *repository.MessagesRepository
}

func NewMessagesHandler(messagesRepository *repository.MessagesRepository) (*MessagesHandler, error) {
	return &MessagesHandler{messagesRepository: messagesRepository}, nil
}

// List godoc
//
//	@Summary		List Sent Messages
//	@Description	list sent messages from server
//	@Tags			Messages
//	@Accept			json
//	@Produce		json
//	@Param			page	query	int	false	"Page number"

// @Success		200	{object}	ListResponse
// @Failure		400	{object}	Error
// @Failure		500	{object}	Error
//
//	@Router			/messages/list [get]
func (s *MessagesHandler) List(w http.ResponseWriter, req *http.Request) {
	pageStr := req.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"ERROR\" : \"%s\"}", err.Error())))
		return
	}

	models, err := s.messagesRepository.ListSentMessages(req.Context(), page)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"ERROR\" : \"%s\"}", err.Error())))
		return
	}

	var messages []*Message

	for _, message := range models {
		messages = append(messages, &Message{
			SentAt:            message.UpdatedAt,
			ExternalMessageID: *message.ExternalMessageID,
		})
	}

	responseData, err := json.Marshal(&messages)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"ERROR\" : \"%s\"}", err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseData))
}
