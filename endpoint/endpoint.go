package endpoint

import (
	"context"
)

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Endpoint[Request any, Response any] func(ctx context.Context, request Request) (response Response, err error)

// Nop is an endpoint that does nothing and returns a nil error.
// Useful for tests.
func Nop[Request any, Response any](context.Context, Request) (Response, error) {
	return *new(Response), nil
}

// Middleware is a chainable behavior modifier for endpoints.
type Middleware func(Endpoint[any, any]) Endpoint[any, any]

// Chain is a helper function for composing middlewares. Requests will
// traverse them in the order they're declared. That is, the first middleware
// is treated as the outermost middleware.
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Endpoint[any, any]) Endpoint[any, any] {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

// Failer may be implemented by Go kit response types that contain business
// logic error details. If Failed returns a non-nil error, the Go kit transport
// layer may interpret this as a business logic error, and may encode it
// differently than a regular, successful response.
//
// It's not necessary for your response types to implement Failer, but it may
// help for more sophisticated use cases. The addsvc example shows how Failer
// should be used by a complete application.
type Failer interface {
	Failed() error
}
