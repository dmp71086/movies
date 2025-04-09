package middleware

import (
	"context"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoverUnaryInterceptor - recover panics
func RecoverUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ interface{}, err error) {
		defer func() {
			if v := recover(); v != nil {
				log.Println(v)

				err = status.Error(codes.Internal, codes.Internal.String()) // return error
			}
		}()

		return handler(ctx, req)
	}
}

func RecoverStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{}, 
		ss grpc.ServerStream, 
		info *grpc.StreamServerInfo, 
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if v := recover(); v != nil {
				log.Println(v)

				err = status.Error(codes.Internal, codes.Internal.String()) // return error
			}
		}()

		return handler(srv, ss)
	}

}
