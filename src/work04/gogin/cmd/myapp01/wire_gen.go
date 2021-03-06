// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"gogin/internal/app/myapp01/conf"
)

// Injectors from wire.go:

// InitializeEvent 声明injector的函数签名
func InitializeEvent(msg string) conf.Event {
	message := conf.NewMessage(msg)
	greeter := conf.NewGreeter(message)
	event := conf.NewEvent(greeter)
	return event
}
