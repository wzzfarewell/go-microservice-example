package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/wzzfarewell/go-microservice-example/internal/util"
	"github.com/wzzfarewell/go-microservice-example/pkg/watermark/endpoint"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

func NewHTTPHandler(eps endpoint.Set) http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").Path("/api/v1/watermark/healthz").Handler(httptransport.NewServer(
		eps.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/api/v1/watermark/documents/{id}/status").Handler(httptransport.NewServer(
		eps.StatusEndpoint,
		decodeHTTPStatusRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/api/v1/watermark/documents").Handler(httptransport.NewServer(
		eps.FindEndpoint,
		decodeHTTPFindRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/api/v1/watermark/documents").Handler(httptransport.NewServer(
		eps.CreateDocumentEndpoint,
		decodeHTTPCreateDocumentRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/api/v1/watermark/watermark").Handler(httptransport.NewServer(
		eps.WatermarkEndpoint,
		decodeHTTPWatermarkRequest,
		encodeResponse,
	))

	return r
}

func routerHandler(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handler.ServeHTTP(w, r)
	}
}

func decodeHTTPFindRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.FindRequest
	if r.ContentLength == 0 {
		logger.Log("Get request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.StatusRequest
	vars := mux.Vars(r)
	req.TicketID = vars["id"]
	return req, nil
}

func decodeHTTPWatermarkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.WatermarkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPCreateDocumentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.CreateDocumentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoint.ServiceStatusRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case util.ErrPathParamNotFound:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
