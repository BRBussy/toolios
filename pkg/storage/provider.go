package storage

import "context"

type Provider interface {
	Store(ctx context.Context, request StoreRequest) (*StoreResponse, error)
}

type StoreRequest struct {
	Data     []byte `validate:"required"`
	MIMEType string `validate:"required"`
	Path     string `validate:"required"`
}

type StoreResponse struct {
}
