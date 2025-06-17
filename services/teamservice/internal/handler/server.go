package handler

import (
	"context"
	"errors"
	teamv1 "github.com/cms-crs/protos/gen/go/team_service"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"userservice/internal/dto"
)

type TeamService interface {
	CreateTeam(
		ctx context.Context,
		req *dto.CreateTeamRequest,
	) (*dto.CreateTeamResponse, error)
	DeleteTeam(
		ctx context.Context,
		ID string,
	) error
	GetTeam(
		ctx context.Context,
		ID string,
	) (*dto.GetTeamResponse, error)
	UpdateTeam(
		ctx context.Context,
		req *dto.UpdateTeamRequest,
	) (*dto.UpdateTeamResponse, error)
	GetUserTeams(
		ctx context.Context,
		UserID string,
	) ([]*dto.GetTeamResponse, error)
	AddUserToTeam(
		ctx context.Context,
		req *dto.AddUserToTeamRequest,
	) (*dto.TeamMember, error)
	RemoveUserFromTeam(
		ctx context.Context,
		req *dto.RemoveUserFromTeamRequest,
	) error
	UpdateUserRole(
		ctx context.Context,
		req *dto.UpdateUserRoleRequest,
	) (*dto.TeamMember, error)
	GetTeamMembers(
		ctx context.Context,
		ID string,
	) ([]*dto.TeamMember, error)
}

type GrpcHandler struct {
	teamv1.UnimplementedTeamServiceServer
	log         *slog.Logger
	teamService TeamService
}

func Register(gRPC *grpc.Server, log *slog.Logger, teamService TeamService) {
	teamv1.RegisterTeamServiceServer(gRPC, &GrpcHandler{
		log:         log,
		teamService: teamService,
	})
}

func (handler *GrpcHandler) CreateTeam(ctx context.Context, request *teamv1.CreateTeamRequest) (*teamv1.Team, error) {
	const op = "gRPC.CreateTeam"

	createTeamRequest := &dto.CreateTeamRequest{
		Name:        request.Name,
		Description: request.Description,
		CreatedBy:   request.CreatedBy,
	}

	createTeamResponse, err := handler.teamService.CreateTeam(ctx, createTeamRequest)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &teamv1.Team{
		Id:          createTeamResponse.ID,
		Name:        createTeamResponse.Name,
		Description: createTeamResponse.Description,
		CreatedAt:   timestamppb.New(createTeamResponse.CreatedAt),
		UpdatedAt:   timestamppb.New(createTeamResponse.UpdatedAt),
	}, nil
}

func (handler *GrpcHandler) DeleteTeam(ctx context.Context, request *teamv1.DeleteTeamRequest) (*empty.Empty, error) {
	const op = "gRPC.DeleteTeam"

	err := handler.teamService.DeleteTeam(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete team")
	}

	return &empty.Empty{}, nil
}

func (handler *GrpcHandler) GetTeam(ctx context.Context, request *teamv1.GetTeamRequest) (*teamv1.Team, error) {
	const op = "gRPC.GetTeam"

	getTeamResponse, err := handler.teamService.GetTeam(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get team")
	}

	return &teamv1.Team{
		Id:          getTeamResponse.ID,
		Name:        getTeamResponse.Name,
		Description: getTeamResponse.Description,
		CreatedAt:   timestamppb.New(getTeamResponse.CreatedAt),
		UpdatedAt:   timestamppb.New(getTeamResponse.UpdatedAt),
	}, nil
}

func (handler *GrpcHandler) UpdateTeam(ctx context.Context, request *teamv1.UpdateTeamRequest) (*teamv1.Team, error) {
	const op = "gRPC.UpdateTeam"

	updateTeamRequest := &dto.UpdateTeamRequest{
		ID:          request.Id,
		Name:        request.Name,
		Description: request.Description,
	}

	updateTeamResponse, err := handler.teamService.UpdateTeam(ctx, updateTeamRequest)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to updated team")
	}

	return &teamv1.Team{
		Id:          updateTeamResponse.ID,
		Name:        updateTeamResponse.Name,
		Description: updateTeamResponse.Description,
		CreatedAt:   timestamppb.New(updateTeamResponse.CreatedAt),
		UpdatedAt:   timestamppb.New(updateTeamResponse.UpdatedAt),
	}, nil
}

func (handler *GrpcHandler) AddUserToTeam(ctx context.Context, request *teamv1.AddUserToTeamRequest) (*teamv1.TeamMember, error) {
	const op = "gRPC.AddUserToTeam"

	addUserToTeamRequest := &dto.AddUserToTeamRequest{
		UserID: request.UserId,
		TeamID: request.TeamId,
		Role:   request.Role,
	}

	addUserToTeamResponse, err := handler.teamService.AddUserToTeam(ctx, addUserToTeamRequest)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.InvalidArgument, "user already in team")
		}
		return nil, status.Error(codes.Internal, "failed to add user to team")
	}

	// todo: implement get user info
	return &teamv1.TeamMember{
		TeamId: addUserToTeamResponse.TeamID,
		UserId: addUserToTeamResponse.UserID,
		Role:   addUserToTeamResponse.Role,
	}, nil
}

func (handler *GrpcHandler) GetUserTeams(ctx context.Context, request *teamv1.GetUserTeamsRequest) (*teamv1.GetUserTeamsResponse, error) {
	getTeamResponses, err := handler.teamService.GetUserTeams(ctx, request.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user teams")
	}

	var getUserTeamsResponse []*teamv1.TeamWithRole
	for _, team := range getTeamResponses {
		getUserTeamsResponse = append(getUserTeamsResponse, &teamv1.TeamWithRole{
			Team: &teamv1.Team{
				Id:          team.ID,
				Name:        team.Name,
				Description: team.Description,
				CreatedAt:   timestamppb.New(team.CreatedAt),
				UpdatedAt:   timestamppb.New(team.UpdatedAt),
			},
		})
	}

	return &teamv1.GetUserTeamsResponse{
		Teams: getUserTeamsResponse,
	}, nil
}

func (handler *GrpcHandler) RemoveUserFromTeam(ctx context.Context, request *teamv1.RemoveUserFromTeamRequest) (*empty.Empty, error) {
	const op = "gRPC.RemoveUserFromTeam"

	removeUserFromTeamRequest := &dto.RemoveUserFromTeamRequest{
		TeamID: request.TeamId,
		UserID: request.UserId,
	}
	err := handler.teamService.RemoveUserFromTeam(ctx, removeUserFromTeamRequest)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to remove user from team")
	}

	return &empty.Empty{}, nil
}

func (handler *GrpcHandler) UpdateUserRole(ctx context.Context, request *teamv1.UpdateUserRoleRequest) (*teamv1.TeamMember, error) {
	updateUserRoleRequest := &dto.UpdateUserRoleRequest{
		TeamID: request.TeamId,
		UserID: request.UserId,
		Role:   request.Role,
	}

	updateUserRoleResponse, err := handler.teamService.UpdateUserRole(ctx, updateUserRoleRequest)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to updated user role")
	}

	return &teamv1.TeamMember{
		UserId: updateUserRoleResponse.UserID,
		TeamId: updateUserRoleResponse.TeamID,
		Role:   updateUserRoleResponse.Role,
	}, nil
}

func (handler *GrpcHandler) GetTeamMembers(ctx context.Context, request *teamv1.GetTeamMembersRequest) (*teamv1.GetTeamMembersResponse, error) {
	teamMembers, err := handler.teamService.GetTeamMembers(ctx, request.TeamId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get team members")
	}

	var getTeamMembersResponse []*teamv1.TeamMember
	for _, teamMember := range teamMembers {
		getTeamMembersResponse = append(getTeamMembersResponse, &teamv1.TeamMember{
			UserId: teamMember.UserID,
			TeamId: teamMember.TeamID,
			Role:   teamMember.Role,
		})
		// add user
	}

	return &teamv1.GetTeamMembersResponse{
		Members: getTeamMembersResponse,
	}, nil
}
