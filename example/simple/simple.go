package main

import (
	"context"
	"log"

	_ "github.com/souma946/gopool"
)

const (
	key = "num"
)

func main() {
	d := pool.NewWorkerPool(1)
	ctx := context.Background()
	ctx = context.WithValue(ctx, key, i)
	d.Execute(ctx, func(ctx context.Context) error {
		arg := ctx.Value(key).(int)
		log.Printf("hello world %d\n", arg)
		return nil
	})
	d.Shutdown(10)
}
