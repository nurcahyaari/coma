package service

import "context"

type ApiKeyServicer interface {
	AuthServicer
	CreateApplicationKey(ctx context.Context) error
}
