package main

import (
	"context"
	"log"
	"net"
	"os"

	"tablelink_project/proto/api"
	"tablelink_project/server/controller"
	mid "tablelink_project/server/middleware"
	"tablelink_project/server/model"
	"tablelink_project/server/repository"
	"tablelink_project/server/service"
	"tablelink_project/server/utils"

	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func main() {
	db := model.BuildDB()
	defer func() {
		if sqlDB, err := db.DB(); err != nil {
			panic(err)
		} else {
			_ = sqlDB.Close()
		}
	}()

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		DB:   0,
	})
	defer redisClient.Close()
	log.Println("Connected to Redis")

	// Test Redis connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	roleValidator := utils.NewRoleRightsValidator(db)
	authController := controller.NewAuthController(userService, redisClient)
	userController := controller.NewUserController(userService, *roleValidator)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			mid.JwtAuthInterceptor(), // Add the JWT Auth middleware
		)),
	)
	api.RegisterAuthServiceServer(grpcServer, authController)
	api.RegisterUserServiceServer(grpcServer, userController)

	log.Println("gRPC server is running at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
