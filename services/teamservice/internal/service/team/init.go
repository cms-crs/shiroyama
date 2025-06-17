package team

import (
	"context"
	"log/slog"
	"userservice/internal/dto"
	"userservice/internal/entity"
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
}

type Service struct {
	log            *slog.Logger
	teamRepository Repository
}

func NewTeamService(log *slog.Logger, teamRepository Repository) *Service {
	return &Service{
		log:            log,
		teamRepository: teamRepository,
	}
}
