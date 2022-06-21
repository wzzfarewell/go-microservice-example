package watermark

import (
	"context"

	"github.com/wzzfarewell/go-microservice-example/internal"
)

type Service interface {
	Find(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error)
	Status(ctx context.Context, ticketID string) (internal.Status, error)
	Watermark(ctx context.Context, ticketID, mark string) (int, error)
	CreateDocument(ctx context.Context, doc *internal.Document) (string, error)
	ServiceStatus(ctx context.Context) (int, error)
}
