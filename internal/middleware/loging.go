package middleware

import (
	"context"

	"movies/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogErrorUnaryInterceptor - log interceptor
func LogErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		id := uuid.NewString()

		logger.Global.Info("request received",
			zap.String("method", info.FullMethod),
			zap.String("requestId", id),
		)

		resp, err = handler(ctx, req)

		logger.Global.Info("request handled",
			zap.String("method", info.FullMethod),
			zap.String("requestId", id),
		)

		if err != nil {
			// 4ХХ -> warn
			// 5ХХ -> Error
			logger.Global.Error(err.Error())

			if status.Code(err) == codes.Internal {
				err = status.Error(codes.Internal, codes.Internal.String())
			}
		}

		return
	}
}

func LoginStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		id := uuid.NewString()

		logger.Global.Info("request received",
			zap.String("method", info.FullMethod),
			zap.String("requestId", id),
		)

		err = handler(srv, ss)

		logger.Global.Info("request handled",
			zap.String("method", info.FullMethod),
			zap.String("requestId", id),
		)

		if err != nil {
			// 4ХХ -> warn
			// 5ХХ -> Error
			logger.Global.Error(err.Error())

			if status.Code(err) == codes.Internal {
				err = status.Error(codes.Internal, codes.Internal.String())
			}
		}

		return
	}
}
