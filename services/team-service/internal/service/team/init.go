package team

import (
	"context"
	"log/slog"
	"taskservice/internal/clients"
	"taskservice/internal/dto"
	"taskservice/internal/entity"
	"taskservice/internal/kafka"
)

type Repository interface {
	CreateTeam(
		ctx context.Context,
		req *entity.Team,
	) (*entity.Team, error)
	DeleteTeam(
		ctx context.Context,
		ID string,
	) error
	GetTeam(
		ctx context.Context,
		ID string,
	) (*entity.Team, error)
	UpdateTeam(
		ctx context.Context,
		team *entity.Team,
	) (*entity.Team, error)
	GetUserTeams(
		ctx context.Context,
		UserID string,
	) ([]*entity.Team, error)
	AddUserToTeam(
		ctx context.Context,
		req *entity.TeamMember,
	) error
	RemoveUserFromTeam(
		ctx context.Context,
		req *dto.RemoveUserFromTeamRequest,
	) error
	UpdateUserRole(
		ctx context.Context,
		req *dto.UpdateUserRoleRequest,
	) error
	GetTeamMembers(
		ctx context.Context,
		ID string,
	) ([]*entity.TeamMember, error)
	DeleteUserFromAllTeams(
		ctx context.Context,
		userID string,
	) (*kafka.TeamDeletionData, error)
	RestoreUserTeams(
		ctx context.Context,
		userID string,
		data *kafka.TeamDeletionData,
	) error
	GetUserTeamMemberships(
		ctx context.Context,
		userID string,
	) (*kafka.TeamDeletionData, error)
}

type Service struct {
	log            *slog.Logger
	teamRepository Repository
	userClient     *clients.UserClient
}

func NewTeamService(log *slog.Logger, teamRepository Repository, userClient *clients.UserClient) *Service {
	return &Service{
		log:            log,
		teamRepository: teamRepository,
		userClient:     userClient,
	}
}
