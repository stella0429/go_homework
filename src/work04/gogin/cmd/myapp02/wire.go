// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"gogin/internal/app/myapp02/biz"
	"gogin/internal/app/myapp02/conf"
	"gogin/internal/app/myapp02/data"
	"gogin/internal/app/myapp02/server"
	"gogin/internal/app/myapp02/service"
)

// initApp
func initApp(*conf.Server, *conf.Data) (*App, error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
