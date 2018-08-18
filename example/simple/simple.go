package main

import (
	"context"
	"log"

	"github.com/souma946/gopool"
)

type (
	sampleKey string
)

const (
	key sampleKey = "num"
)

func main() {
	d := gopool.NewWorkerPool(1)
	ctx := context.Background()
	ctx = context.WithValue(ctx, key, 100)
	d.Execute(ctx, func(ctx context.Context) error {
		arg := ctx.Value(key).(int)
		log.Printf("hello world %d\n", arg)
		return nil
	})
	d.Shutdown(10)
}
