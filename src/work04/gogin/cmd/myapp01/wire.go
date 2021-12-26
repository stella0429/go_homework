// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"gogin/internal/app/myapp01/conf"
)

// InitializeEvent 声明injector的函数签名
func InitializeEvent(msg string) conf.Event {
	wire.Build(conf.NewEvent, conf.NewGreeter, conf.NewMessage)
	return conf.Event{} //返回值没有实际意义，只需符合函数签名即可
}
