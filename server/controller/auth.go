package controller

import (
	"context"
	"errors"
	"os"
	"strconv"
	pb "tablelink_project/proto/api"
	"tablelink_project/server/service"
	"time"

	"github.com/go-redis/redis/v8"
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

	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		response.Status = false
		response.Message = "Failed to parse token lifespan: " + err.Error()
		return response, errors.New("failed to save token in Redis")
	}

	err = ac.redisClient.Set(ctx, token, req.Email, time.Duration(token_lifespan)).Err()
	if err != nil {
		response.Status = false
		response.Message = "Failed to save token in Redis: " + err.Error()
		return response, errors.New("failed to save token in Redis")
	}

	return &pb.LoginResponse{
		Status:      true,
		Message:     "Login successful",
		AccessToken: token,
	}, nil
}
