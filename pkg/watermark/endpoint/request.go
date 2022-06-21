package endpoint

import "github.com/wzzfarewell/go-microservice-example/internal"

type FindRequest struct {
	Filters []internal.Filter `json:"filters,omitempty"`
}

type StatusRequest struct {
	TicketID string `json:"ticket_id"`
}

type WatermarkRequest struct {
	TicketID string `json:"ticket_id"`
	Mark     string `json:"mark"`
}

type CreateDocumentRequest struct {
	Document *internal.Document `json:"document"`
}

type ServiceStatusRequest struct{}
