package main

import (
	"context"
)

type LocalMap struct {
	data map[string][]byte
}

func (l *LocalMap) Set(_ context.Context, key string, value []byte) (err error) {
	if l.data == nil {
		l.data = map[string][]byte{}
	}
	l.data[key] = value
	return nil
}

func (l *LocalMap) Get(_ context.Context, key string) (value []byte, err error) {
	if l.data == nil {
		l.data = map[string][]byte{}
	}
	value = l.data[key]
	return value, nil
}

func (l *LocalMap) Delete(_ context.Context, key string) (err error) {
	if l.data == nil {
		return nil
	}
	delete(l.data, key)
	return nil
}
