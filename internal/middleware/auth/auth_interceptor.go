package auth

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *JWTManager) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func (interceptor *JWTManager) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		err := interceptor.authorize(stream.Context())
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *JWTManager) authorize(ctx context.Context) error {
	// accessibleRoles, ok := interceptor.accessibleRoles[method]
	// if !ok {
	// 	// everyone can access
	// 	return nil
	// }

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) != 3 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	_, err := interceptor.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	// for _, role := range accessibleRoles {
	// 	if role == claims.Role {
	// 		return nil
	// 	}
	// }

	// return status.Error(codes.PermissionDenied, "no permission to access this RPC")

	return nil
}
