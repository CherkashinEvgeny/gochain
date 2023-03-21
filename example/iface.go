package main

import "context"

type Map interface {
	Set(ctx context.Context, key string, value []byte) (err error)
	Get(ctx context.Context, key string) (value []byte, err error)
	Delete(ctx context.Context, key string) (err error)
}
