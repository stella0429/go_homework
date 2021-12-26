package data

import (
	"github.com/google/wire"
	"gogin/internal/app/myapp02/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data) (*Data, error) {
	return &Data{}, nil
}
