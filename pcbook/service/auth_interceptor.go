package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(
	jwtManager *JWTManager,
	accessibleRoles map[string][]string,
) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Panicln("--> unary interceptor: ", info.FullMethod)

		// TODO: implement authorization

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Panicln("--> stream interceptor: ", info.FullMethod)

		// TODO: implement authorization

		return handler(srv, stream)
	}
}
