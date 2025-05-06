package middleware

import (
	"context"
	"net/http"
	"tablelink_project/server/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func JwtAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(http.StatusUnauthorized, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Errorf(http.StatusUnauthorized, "missing authorization token")
		}
		token := tokens[0]

		if info.FullMethod != "/auth/login" {
			userID, err := utils.ValidateToken(token)
			if err != nil {
				return nil, status.Errorf(http.StatusUnauthorized, "invalid token: %v", err)
			}

			ctx = context.WithValue(ctx, utils.UserCtxKey, userID)
		}

		return handler(ctx, req)
	}
}
