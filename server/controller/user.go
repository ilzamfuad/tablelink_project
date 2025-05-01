package controller

import (
	"context"
	"fmt"
	"tablelink_project/config"
	pb "tablelink_project/proto/api"
	"tablelink_project/server/model"
	"tablelink_project/server/service"
	"tablelink_project/server/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserController struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	err := uc.roleValidate(ctx)
	response := &pb.GetAllUsersResponse{}

	if err != nil {
		response.Status = false
		response.Message = fmt.Sprintf("error: %v", err.Error())
		return response, status.Errorf(codes.PermissionDenied, "access denied: %v", err.Error())
	}

	users, err := uc.userService.GetAllUsers()
	if err != nil {
		response.Status = false
		response.Message = fmt.Sprintf("error: %v", err.Error())
		return response, status.Errorf(codes.Internal, "error: %v", err.Error())
	}

	for _, user := range users {
		response.Message = "success"
		response.Status = true
		response.Data = append(response.Data, &pb.User{
			UserId:     uint32(user.ID),
			Email:      user.Email,
			RoleId:     uint32(user.Role.ID),
			RoleName:   user.Role.Name,
			LastAccess: user.LastAccess.Format("2006-01-02 15:04:05"),
		})
	}
	return response, nil

}

func (uc *UserController) CreateUserRequest(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	err := uc.roleValidate(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: %v", err.Error())
	}
	response := &pb.CreateUserResponse{}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   uint(req.RoleId),
	}
	err = uc.userService.CreateUser(user)
	if err != nil {
		response.Status = false
		response.Message = fmt.Sprintf("error: %v", err.Error())
		return response, status.Errorf(codes.Internal, "error: %v", err.Error())
	}

	response.Status = true
	response.Message = "success"
	return response, nil
}

func (uc *UserController) UpdateUserRequest(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := uc.roleValidate(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: %v", err.Error())
	}

	response := &pb.UpdateUserResponse{}
	userID, ok := ctx.Value(utils.UserCtxKey).(uint)
	if !ok {
		response.Status = false
		response.Message = "user not authenticated"
		return response, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	user, err := uc.userService.GetUserByID(int(userID))
	if err != nil {
		response.Status = false
		response.Message = fmt.Sprintf("error: %v", err.Error())
		return response, status.Errorf(codes.Internal, "error: %v", err.Error())
	}

	user.Name = req.Name

	err = uc.userService.UpdateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err.Error())
	}

	return &pb.UpdateUserResponse{
		Status:  true,
		Message: "success",
	}, nil
}

func (uc *UserController) DeleteUserRequest(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := uc.roleValidate(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: %v", err.Error())
	}

	response := &pb.DeleteUserResponse{}
	err = uc.userService.DeleteUser(int(req.UserId))
	if err != nil {
		response.Status = false
		response.Message = fmt.Sprintf("error: %v", err.Error())
		return response, status.Errorf(codes.Internal, "error: %v", err.Error())
	}
	return &pb.DeleteUserResponse{
		Status:  true,
		Message: "success",
	}, nil
}

func (uc *UserController) roleValidate(ctx context.Context) error {
	userID, ok := ctx.Value(utils.UserCtxKey).(uint)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	sections := md.Get("X-Link-Service")
	if len(sections) == 0 {
		return status.Errorf(codes.Unauthenticated, "missing section in metadata")
	}
	section := sections[0]

	route, _ := grpc.Method(ctx)

	restMapping, exists := config.GrpcToRestfulMapping[route]
	if !exists {
		return status.Errorf(codes.InvalidArgument, "invalid route: %s", route)
	}

	err := uc.userService.ValidateRoleRights(userID, section, restMapping.Route, restMapping.Method)
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "access denied: %v", err.Error())
	}

	return nil
}
