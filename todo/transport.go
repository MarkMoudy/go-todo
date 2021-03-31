package todo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeRoutes mounts all routes for this service onto the router.
func MakeRoutes(r *mux.Router, s Service, logger log.Logger) {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	for _, rt := range httpRoutes {
		route := r.Methods(rt.method).Path(rt.pattern)

		h := kithttp.NewServer(
			rt.getEndpoint(s),
			rt.decReq,
			encodeResponse,
			opts...,
		)
		route.Handler(h)
	}
}

var httpRoutes = []struct {
	pattern     string
	method      string
	getEndpoint func(s Service) endpoint.Endpoint
	decReq      kithttp.DecodeRequestFunc
}{
	{
		pattern:     "/v1/todos/{id:[0-9]+}/status",
		method:      http.MethodPut,
		getEndpoint: func(s Service) endpoint.Endpoint { return makeSetCompletedStatusEndpoint(s) },
		decReq:      decodeSetCompletedStatusRequest,
	},
	{
		pattern:     "/v1/todos/{id:[0-9]+}",
		method:      http.MethodGet,
		getEndpoint: func(s Service) endpoint.Endpoint { return makeTodoEndpoint(s) },
		decReq:      decodeIDRequest,
	},
	{
		pattern:     "/v1/todos/{id:[0-9]+}",
		method:      http.MethodPut,
		getEndpoint: func(s Service) endpoint.Endpoint { return makeUpdateTodoEndpoint(s) },
		decReq:      decodeUpdateTodoRequest,
	},
	{
		pattern:     "/v1/todos",
		method:      http.MethodGet,
		getEndpoint: func(s Service) endpoint.Endpoint { return makeTodosEndpoint(s) },
		decReq:      kithttp.NopRequestDecoder,
	},
	{
		pattern:     "/v1/todos",
		method:      http.MethodPost,
		getEndpoint: func(s Service) endpoint.Endpoint { return makeCreateTodoEndpoint(s) },
		decReq:      decodeCreateTodoRequest,
	},
}

func decodeIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	return idRequest{id: id}, nil
}

func parseID(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return 0, errBadRouting
	}

	return strconv.Atoi(id)
}

func decodeCreateTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	dec := json.NewDecoder(r.Body)

	var req createTodoRequest
	if err := dec.Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	dec := json.NewDecoder(r.Body)

	var req updateTodoRequest
	if err := dec.Decode(&req); err != nil {
		return nil, err
	}

	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	req.ID = id

	return req, nil
}

func decodeSetCompletedStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	dec := json.NewDecoder(r.Body)

	var req setCompletedStatusRequest
	if err := dec.Decode(&req); err != nil {
		return nil, err
	}

	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	req.ID = id

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError

	switch {
	case errors.Is(err, ErrInvalidID):
		code = http.StatusBadRequest
	case errors.Is(err, ErrTodoNotFound):
		code = http.StatusNotFound
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
