package main

import (
	"context"
	"fmt"
)

func main() {
	chain := MapChain{}
	chain.Register(10, func(m Map) Map {
		return &LocalMap{}
	})
	chain.Register(8, func(m Map) Map {
		return Interceptor{m}
	})
	m := chain.Instance()
	_ = m.Set(context.Background(), "test", []byte{0})
	fmt.Println(m.Get(context.Background(), "test"))
	_ = m.Set(context.Background(), "lucky", []byte{0})
	fmt.Println(m.Get(context.Background(), "lucky"))
}

type Interceptor struct {
	impl Map
}

func (i Interceptor) Set(ctx context.Context, key string, value []byte) (err error) {
	if key == "lucky" {
		return nil
	}
	return i.impl.Set(ctx, key, value)
}

func (i Interceptor) Get(ctx context.Context, key string) (value []byte, err error) {
	return i.impl.Get(ctx, key)
}

func (i Interceptor) Delete(ctx context.Context, key string) (err error) {
	return i.impl.Delete(ctx, key)
}
