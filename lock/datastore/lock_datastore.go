package datastore

import (
	"context"
)

type DataStore interface {
	Create(ctx context.Context, collection string, data interface{}) error
	Read(ctx context.Context, collection string, filter interface{}, result interface{}) error
	Delete(ctx context.Context, collection string, filter interface{}) error
}
