package transport

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var wsSubscriptionErrorCtxKey = &wsSubscriptionErrorContextKey{"subscription-error"}

type wsSubscriptionErrorContextKey struct {
	name string
}

type subscriptionError struct {
	errs []*gqlerror.Error
}

// AddSubscriptionError is used to let websocket return an error message after subscription resolver returns a channel.
// for example:
//
//	func (r *subscriptionResolver) Method(ctx context.Context) (<-chan *model.Message, error) {
//		ch := make(chan *model.Message)
//		go func() {
//	     defer func() {
//				close(ch)
//	     }
//			// some kind of block processing (e.g.: gRPC client streaming)
//			stream, err := gRPCClientStreamRequest(ctx)
//			if err != nil {
//				   transport.AddSubscriptionError(ctx, err)
//	            return // must return and close channel so websocket can send error back
//	     }
//			for {
//				m, err := stream.Recv()
//				if err == io.EOF {
//					return
//				}
//				if err != nil {
//				   transport.AddSubscriptionError(ctx, err)
//	            return // must return and close channel so websocket can send error back
//				}
//				ch <- m
//			}
//		}()
//
//		return ch, nil
//	}
//
// see https://github.com/99designs/gqlgen/pull/2506 for more details
func AddSubscriptionError(ctx context.Context, err *gqlerror.Error) {
	subscriptionErrStruct := getSubscriptionErrorStruct(ctx)
	subscriptionErrStruct.errs = append(subscriptionErrStruct.errs, err)
}

func withSubscriptionErrorContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, wsSubscriptionErrorCtxKey, &subscriptionError{})
}

func getSubscriptionErrorStruct(ctx context.Context) *subscriptionError {
	v, _ := ctx.Value(wsSubscriptionErrorCtxKey).(*subscriptionError)
	return v
}

func getSubscriptionError(ctx context.Context) []*gqlerror.Error {
	return getSubscriptionErrorStruct(ctx).errs
}
