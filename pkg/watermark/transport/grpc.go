package transport

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/wzzfarewell/go-microservice-example/api/v1/pb/watermark"
	"github.com/wzzfarewell/go-microservice-example/internal"
	"github.com/wzzfarewell/go-microservice-example/pkg/watermark/endpoint"
)

type grpcServer struct {
	watermark.UnimplementedWatermarkServer
	find           grpc.Handler
	status         grpc.Handler
	serviceStatus  grpc.Handler
	createDocument grpc.Handler
	watermark      grpc.Handler
}

func NewGRPCServer(ep endpoint.Set) watermark.WatermarkServer {
	return &grpcServer{
		find:           grpc.NewServer(ep.FindEndpoint, decodeGRPCGetRequest, encodeGRPCGetResponse),
		status:         grpc.NewServer(ep.StatusEndpoint, decodeGRPCStatusRequest, encodeGRPCStatusResponse),
		serviceStatus:  grpc.NewServer(ep.ServiceStatusEndpoint, decodeGRPCServiceStatusRequest, encodeGRPCServiceStatusResponse),
		createDocument: grpc.NewServer(ep.CreateDocumentEndpoint, decodeGRPCCreateDocumentRequest, encodeGRPCCreateDocumentResponse),
		watermark:      grpc.NewServer(ep.WatermarkEndpoint, decodeGRPCWatermarkRequest, encodeGRPCWatermarkResponse),
	}
}

func (s *grpcServer) Find(ctx context.Context, request *watermark.FindRequest) (*watermark.FindReply, error) {
	_, reply, err := s.find.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply.(*watermark.FindReply), nil
}

func (s *grpcServer) Watermark(ctx context.Context, request *watermark.WatermarkRequest) (*watermark.WatermarkReply, error) {
	_, reply, err := s.watermark.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply.(*watermark.WatermarkReply), nil
}

func (s *grpcServer) Status(ctx context.Context, request *watermark.StatusRequest) (*watermark.StatusReply, error) {
	_, reply, err := s.status.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply.(*watermark.StatusReply), nil
}

func (s *grpcServer) CreateDocument(ctx context.Context, request *watermark.CreateDocumentRequest) (*watermark.CreateDocumentReply, error) {
	_, reply, err := s.createDocument.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply.(*watermark.CreateDocumentReply), nil
}

func (s *grpcServer) ServiceStatus(ctx context.Context, request *watermark.ServiceStatusRequest) (*watermark.ServiceStatusReply, error) {
	_, reply, err := s.serviceStatus.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply.(*watermark.ServiceStatusReply), nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.FindRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoint.FindRequest{Filters: filters}, nil
}

func decodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.StatusRequest)
	return endpoint.StatusRequest{TicketID: req.TicketID}, nil
}

func decodeGRPCWatermarkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.WatermarkRequest)
	return endpoint.WatermarkRequest{TicketID: req.TicketID, Mark: req.Mark}, nil
}

func decodeGRPCCreateDocumentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.CreateDocumentRequest)
	doc := &internal.Document{
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return endpoint.CreateDocumentRequest{Document: doc}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoint.ServiceStatusRequest{}, nil
}

func encodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.FindReply)
	var docs []internal.Document
	for _, d := range reply.Documents {
		doc := internal.Document{
			Content:   d.Content,
			Title:     d.Title,
			Author:    d.Author,
			Topic:     d.Topic,
			Watermark: d.Watermark,
		}
		docs = append(docs, doc)
	}
	return endpoint.FindResponse{Documents: docs, Err: reply.Err}, nil
}

func encodeGRPCStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.StatusReply)
	return endpoint.StatusResponse{Status: internal.Status(reply.Status), Err: reply.Err}, nil
}

func encodeGRPCWatermarkResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.WatermarkReply)
	return endpoint.WatermarkResponse{Code: int(reply.Code), Err: reply.Err}, nil
}

func encodeGRPCCreateDocumentResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.CreateDocumentReply)
	return endpoint.CreateDocumentResponse{TicketID: reply.TicketID, Err: reply.Err}, nil
}

func encodeGRPCServiceStatusResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.ServiceStatusReply)
	return endpoint.ServiceStatusResponse{Code: int(reply.Code), Err: reply.Err}, nil
}
