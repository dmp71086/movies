package middleware

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var unaryLimiter = rate.NewLimiter(100, 130)

// LimiterInterceptor - convert any arror to rpc error
func LimiterUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if !unaryLimiter.Allow() {
			return nil, status.Error(codes.ResourceExhausted, "too many request")
		}

		return handler(ctx, req)
	}
}

var streamLimiter = rate.NewLimiter(10, 13)

func LimiterStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		if !streamLimiter.Allow() {
			return status.Error(codes.ResourceExhausted, "too many request")
		}

		return handler(srv, ss)
	}

}
