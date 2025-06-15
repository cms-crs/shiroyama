package handler

import (
	"context"
	teamv1 "github.com/cms-crs/protos/go/proto/team"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
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
	GetUserTeams(
		ctx context.Context,
		UserID string,
	) (*dto.GetUserTeamsResponse, error)
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
		req *dto.UpdateUserRole,
	) (*dto.TeamMember, error)
	GetTeamMembers(
		ctx context.Context,
		ID string,
	) (*dto.GetTeamMembersResponse, error)
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
	return nil, nil
}

func (handler *GrpcHandler) DeleteTeam(ctx context.Context, request *teamv1.DeleteTeamRequest) (*empty.Empty, error) {
	return nil, nil
}

func (handler *GrpcHandler) GetTeam(ctx context.Context, request *teamv1.GetTeamRequest) (*teamv1.Team, error) {
	return nil, nil
}

func (handler *GrpcHandler) GetUserTeams(ctx context.Context, request *teamv1.GetUserTeamsRequest) (*teamv1.GetUserTeamsResponse, error) {
	return nil, nil
}

func (handler *GrpcHandler) AddUserToTeam(ctx context.Context, request *teamv1.AddUserToTeamRequest) (*teamv1.TeamMember, error) {
	return nil, nil
}

func (handler *GrpcHandler) RemoveUserFromTeam(ctx context.Context, request *teamv1.RemoveUserFromTeamRequest) (*empty.Empty, error) {
	return nil, nil
}

func (handler *GrpcHandler) UpdateUserRole(ctx context.Context, request *teamv1.UpdateUserRoleRequest) (*teamv1.TeamMember, error) {
	return nil, nil
}

func (handler *GrpcHandler) GetTeamMembers(ctx context.Context, request *teamv1.GetTeamMembersRequest) (*teamv1.GetTeamMembersResponse, error) {
	return nil, nil
}
