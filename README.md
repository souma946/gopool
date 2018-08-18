# gopool
[![Build Status](https://travis-ci.org/souma946/gopool.svg?branch=master)](https://travis-ci.org/souma946/gopool) 

gopool is useful when you want to control the number of concurrent executions


## Install
```
go get -u github.dom/souma946/gopool
```

## Usage
```
import "github.com/souma946/gopool"

d := gopool.NewWorkerPool(1)
ctx := context.Background()
ctx = context.WithValue(ctx, key, 100)
d.Execute(ctx, func(ctx context.Context) error {
	arg := ctx.Value(key).(int)
	log.Printf("hello world %d\n", arg)
	return nil
})
d.Shutdown(10)
```