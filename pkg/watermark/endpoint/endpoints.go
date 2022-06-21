package endpoint

import (
	"context"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/wzzfarewell/go-microservice-example/pkg/watermark"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

type Set struct {
	FindEndpoint           endpoint.Endpoint
	CreateDocumentEndpoint endpoint.Endpoint
	StatusEndpoint         endpoint.Endpoint
	ServiceStatusEndpoint  endpoint.Endpoint
	WatermarkEndpoint      endpoint.Endpoint
}

func NewEndpointSet(svc watermark.Service) Set {
	return Set{
		FindEndpoint:           MakeFindEndpoint(svc),
		CreateDocumentEndpoint: MakeCreateDocumentEndpoint(svc),
		StatusEndpoint:         MakeStatusEndpoint(svc),
		ServiceStatusEndpoint:  MakeServiceStatusEndpoint(svc),
		WatermarkEndpoint:      MakeWatermarkEndpoint(svc),
	}
}

func MakeFindEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindRequest)
		docs, err := svc.Find(ctx, req.Filters...)
		if err != nil {
			return FindResponse{docs, err.Error()}, nil
		}
		return FindResponse{docs, ""}, nil
	}
}

func MakeStatusEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StatusRequest)
		status, err := svc.Status(ctx, req.TicketID)
		if err != nil {
			return StatusResponse{Status: status, Err: err.Error()}, nil
		}
		return StatusResponse{Status: status, Err: ""}, nil
	}
}

func MakeCreateDocumentEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateDocumentRequest)
		ticketID, err := svc.CreateDocument(ctx, req.Document)
		if err != nil {
			return CreateDocumentResponse{TicketID: ticketID, Err: err.Error()}, nil
		}
		return CreateDocumentResponse{TicketID: ticketID, Err: ""}, nil
	}
}

func MakeWatermarkEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WatermarkRequest)
		code, err := svc.Watermark(ctx, req.TicketID, req.Mark)
		if err != nil {
			return WatermarkResponse{Code: code, Err: err.Error()}, nil
		}
		return WatermarkResponse{Code: code, Err: ""}, nil
	}
}

func MakeServiceStatusEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

// func (s *Set) Find(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
// 	resp, err := s.FindEndpoint(ctx, FindRequest{Filters: filters})
// 	if err != nil {
// 		return []internal.Document{}, err
// 	}
// 	findResp := resp.(FindResponse)
// 	if findResp.Err != "" {
// 		return []internal.Document{}, errors.New(findResp.Err)
// 	}
// 	return findResp.Documents, nil
// }

// func (s *Set) ServiceStatus(ctx context.Context) (int, error) {
// 	resp, err := s.ServiceStatusEndpoint(ctx, ServiceStatusRequest{})
// 	svcStatusResp := resp.(ServiceStatusResponse)
// 	if err != nil {
// 		return svcStatusResp.Code, err
// 	}
// 	if svcStatusResp.Err != "" {
// 		return svcStatusResp.Code, errors.New(svcStatusResp.Err)
// 	}
// 	return svcStatusResp.Code, nil
// }

// func (s *Set) CreateDocument(ctx context.Context, doc *internal.Document) (string, error) {
// 	resp, err := s.CreateDocumentEndpoint(ctx, CreateDocumentRequest{Document: doc})
// 	if err != nil {
// 		return "", err
// 	}
// 	adResp := resp.(CreateDocumentResponse)
// 	if adResp.Err != "" {
// 		return "", errors.New(adResp.Err)
// 	}
// 	return adResp.TicketID, nil
// }

// func (s *Set) Status(ctx context.Context, ticketID string) (internal.Status, error) {
// 	resp, err := s.StatusEndpoint(ctx, StatusRequest{TicketID: ticketID})
// 	if err != nil {
// 		return internal.Failed, err
// 	}
// 	stsResp := resp.(StatusResponse)
// 	if stsResp.Err != "" {
// 		return internal.Failed, errors.New(stsResp.Err)
// 	}
// 	return stsResp.Status, nil
// }

// func (s *Set) Watermark(ctx context.Context, ticketID, mark string) (int, error) {
// 	resp, err := s.WatermarkEndpoint(ctx, WatermarkRequest{TicketID: ticketID, Mark: mark})
// 	wmResp := resp.(WatermarkResponse)
// 	if err != nil {
// 		return wmResp.Code, err
// 	}
// 	if wmResp.Err != "" {
// 		return wmResp.Code, errors.New(wmResp.Err)
// 	}
// 	return wmResp.Code, nil
// }
