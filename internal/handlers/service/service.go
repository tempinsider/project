package service

import (
	"encoding/json"
	"fmt"
	"insider-mert/internal/workers"
	"net/http"
)

type ServiceHandler struct {
	worker *workers.Worker
}

func NewServiceHandler(worker *workers.Worker) (*ServiceHandler, error) {
	return &ServiceHandler{worker: worker}, nil
}

// Toggle godoc
//
//	@Summary		toggle worker status
//	@Description	toggle worker status
//	@Tags			Worker Service
//	@Accept			json
//	@Produce		json
//	@Param			request	body	ToggleRequest	true	"Toggle Request"
//
// @Success		200	{object}	ToggleResponse
// @Failure		500	{object}	Error
// @Router			/services/toggle [post]
func (s *ServiceHandler) Toggle(w http.ResponseWriter, req *http.Request) {
	isWorking := s.worker.IsWorking()
	if isWorking {
		s.worker.Stop()
	} else {
		s.worker.Start()
	}

	response := ToggleResponse{
		ServiceStatus: ServiceStatusWorking,
	}
	if isWorking {
		response.ServiceStatus = ServiceStatusStopped
	}

	responseData, err := json.Marshal(&response)
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
