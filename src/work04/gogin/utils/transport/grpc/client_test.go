package grpc

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gogin/utils/log"
	"gogin/utils/middleware"
	"gogin/utils/registry"
	"google.golang.org/grpc"
)

func TestWithEndpoint(t *testing.T) {
	o := &clientOptions{}
	v := "abc"
	WithEndpoint(v)(o)
	assert.Equal(t, v, o.endpoint)
}

func TestWithTimeout(t *testing.T) {
	o := &clientOptions{}
	v := time.Duration(123)
	WithTimeout(v)(o)
	assert.Equal(t, v, o.timeout)
}

func TestWithMiddleware(t *testing.T) {
	o := &clientOptions{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	WithMiddleware(v...)(o)
	assert.Equal(t, v, o.middleware)
}

type mockRegistry struct{}

func (m *mockRegistry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (m *mockRegistry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return nil, nil
}

func TestWithDiscovery(t *testing.T) {
	o := &clientOptions{}
	v := &mockRegistry{}
	WithDiscovery(v)(o)
	assert.Equal(t, v, o.discovery)
}

func TestWithTLSConfig(t *testing.T) {
	o := &clientOptions{}
	v := &tls.Config{}
	WithTLSConfig(v)(o)
	assert.Equal(t, v, o.tlsConf)
}

func TestWithLogger(t *testing.T) {
	o := &clientOptions{}
	v := log.DefaultLogger
	WithLogger(v)(o)
	assert.Equal(t, v, o.logger)
}

func EmptyMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			return handler(ctx, req)
		}
	}
}

func TestUnaryClientInterceptor(t *testing.T) {
	f := unaryClientInterceptor([]middleware.Middleware{EmptyMiddleware()}, time.Duration(100), nil)
	req := &struct{}{}
	resp := &struct{}{}

	err := f(context.TODO(), "hello", req, resp, &grpc.ClientConn{},
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		})
	assert.NoError(t, err)
}

func TestWithUnaryInterceptor(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.UnaryClientInterceptor{
		func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return nil
		},
		func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return nil
		},
	}
	WithUnaryInterceptor(v...)(o)
	assert.Equal(t, v, o.ints)
}

func TestWithOptions(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithOptions(v...)(o)
	assert.Equal(t, v, o.grpcOpts)
}

func TestDial(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithOptions(v...)(o)
	assert.Equal(t, v, o.grpcOpts)
}

func TestDialConn(t *testing.T) {
	_, err := dial(
		context.Background(),
		true,
		WithDiscovery(&mockRegistry{}),
		WithTimeout(10*time.Second),
		WithLogger(log.DefaultLogger),
		WithEndpoint("abc"),
		WithMiddleware(EmptyMiddleware()),
	)
	if err != nil {
		t.Error(err)
	}
}
