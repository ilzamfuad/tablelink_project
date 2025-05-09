package controller

import (
	"context"
	"errors"
	pb "tablelink_project/proto/api"
	"tablelink_project/server/service"
	"tablelink_project/server/utils"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer
	userService service.UserService
	redisClient *redis.Client
}

func NewAuthController(userService service.UserService, redis *redis.Client) *AuthController {
	return &AuthController{
		userService: userService,
		redisClient: redis,
	}
}

func (ac *AuthController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	response := &pb.LoginResponse{}
	existingToken, err := ac.redisClient.Get(ctx, req.Email).Result()
	if err == nil {
		response.Status = true
		response.Message = "Login successful (from cache)"
		response.AccessToken = existingToken
		return response, nil
	}

	token, err := ac.userService.LoginCheck(req.Email, req.Password)
	if err != nil {
		response.Status = false
		response.Message = "Login failed: " + err.Error()
		return response, errors.New("username or password is incorrect")
	}

	err = ac.redisClient.Set(ctx, req.Email, token, utils.AccessTokenExpiredTime).Err()
	if err != nil {
		response.Status = false
		response.Message = "Failed to save token in Redis: " + err.Error()
		return response, errors.New("failed to save token in Redis")
	}

	return &pb.LoginResponse{
		Status:       true,
		Message:      "Login successful",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (ac *AuthController) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	response := &pb.RefreshTokenResponse{}

	// Validate the refresh token
	token, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		response.Status = false
		response.Message = "Invalid refresh token"
		return response, status.Errorf(codes.Unauthenticated, "invalid refresh token: %v", err)
	}

	// Generate a new access token
	accessToken, err := utils.GenerateToken(token.UserID)
	if err != nil {
		response.Status = false
		response.Message = "Failed to generate access token"
		return response, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	response.Status = true
	response.Message = "Token refreshed successfully"
	response.AccessToken = accessToken.AccessToken
	response.RefreshToken = accessToken.RefreshToken
	return response, nil
}
