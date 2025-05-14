package service

type ToggleRequest struct {
}

type ServiceStatus string

const (
	ServiceStatusStopped ServiceStatus = "STOPPED"
	ServiceStatusWorking ServiceStatus = "WORKING"
)

type Error struct {
	Error error `json:"ERROR"`
}

type ToggleResponse struct {
	ServiceStatus ServiceStatus `json:"SERVICE_STATUS"`
}
