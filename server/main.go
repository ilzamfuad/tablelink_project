package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"tablelink_project/config"
	"tablelink_project/proto/api"
	"tablelink_project/server/controller"
	mid "tablelink_project/server/middleware"
	"tablelink_project/server/repository"
	"tablelink_project/server/service"

	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := config.BuildDB()
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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authController := controller.NewAuthController(userService, redisClient)
	userController := controller.NewUserController(userService)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			mid.JwtAuthInterceptor(),
		)),
	)
	api.RegisterAuthServiceServer(grpcServer, authController)
	api.RegisterUserServiceServer(grpcServer, userController)

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = api.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, ":50051", opts)
	if err != nil {
		log.Fatalf("failed to register AuthService handler: %v", err)
	}

	err = api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, ":50051", opts)
	if err != nil {
		log.Fatalf("failed to register UserService handler: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: mux,
	}

	// Start both servers
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Println("gRPC server is running at :", os.Getenv("GRPC_PORT"))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	log.Println("HTTP server is running at :", os.Getenv("PORT"))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
