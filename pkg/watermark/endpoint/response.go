package endpoint

import "github.com/wzzfarewell/go-microservice-example/internal"

type FindResponse struct {
	Documents []internal.Document `json:"documents"`
	Err       string              `json:"err,omitempty"`
}

type StatusResponse struct {
	Status internal.Status `json:"status"`
	Err    string          `json:"err,omitempty"`
}

type WatermarkResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

type CreateDocumentResponse struct {
	TicketID string `json:"ticket_id"`
	Err      string `json:"err,omitempty"`
}

type ServiceStatusResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}
