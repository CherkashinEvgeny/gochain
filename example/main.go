package main

import (
	"context"
	"fmt"

	"github.com/CherkashinEvgeny/gochain/chain"
)

func main() {
	mapChain := MapChain{}
	mapChain.Register(chain.Impl, func(m Map) Map {
		return &LocalMap{}
	})
	mapChain.Register(chain.Normal, func(m Map) Map {
		return Interceptor{m, chain.Normal}
	})
	mapChain.Register(chain.VeryEarly, func(m Map) Map {
		return Interceptor{m, chain.VeryEarly}
	})
	m := mapChain.Instance()
	_ = m.Set(context.Background(), "test", []byte{0})
	fmt.Println(m.Get(context.Background(), "test"))
	_ = m.Set(context.Background(), "lucky", []byte{0})
	fmt.Println(m.Get(context.Background(), "lucky"))
}

type Interceptor struct {
	impl  Map
	order chain.Order
}

func (i Interceptor) Set(ctx context.Context, key string, value []byte) (err error) {
	fmt.Println("Set: ", i.order)
	if key == "lucky" {
		return nil
	}
	return i.impl.Set(ctx, key, value)
}

func (i Interceptor) Get(ctx context.Context, key string) (value []byte, err error) {
	fmt.Println("Get: ", i.order)
	return i.impl.Get(ctx, key)
}

func (i Interceptor) Delete(ctx context.Context, key string) (err error) {
	fmt.Println("Delete: ", i.order)
	return i.impl.Delete(ctx, key)
}
