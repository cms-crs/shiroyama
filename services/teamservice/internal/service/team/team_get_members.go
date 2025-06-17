package team

import (
	"context"
	"userservice/internal/dto"
)

func (service *Service) GetTeamMembers(ctx context.Context, ID string) ([]*dto.TeamMember, error) {
	const op = "TeamService.GetTeamMembers"

	teamMembers, err := service.teamRepository.GetTeamMembers(ctx, ID)

	if err != nil {
		return nil, err
	}

	var teamMembersDTO []*dto.TeamMember
	for _, teamMember := range teamMembers {
		teamMembersDTO = append(teamMembersDTO, &dto.TeamMember{
			UserID: teamMember.UserID,
			TeamID: teamMember.TeamID,
			Role:   teamMember.Role,
		})
	}

	return teamMembersDTO, nil
}
