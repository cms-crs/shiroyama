package team

import (
	"context"
	"fmt"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"userservice/internal/dto"
)

func (service *Service) GetTeamMembers(ctx context.Context, ID string) ([]*dto.TeamMember, error) {
	const op = "TeamService.GetTeamMembers"

	teamMembers, err := service.teamRepository.GetTeamMembers(ctx, ID)

	if err != nil {
		return nil, err
	}

	var userIDs []string
	for _, member := range teamMembers {
		userIDs = append(userIDs, member.UserID)
	}

	users, err := service.userClient.GetUsersByIds(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	userMap := make(map[string]*userv1.User)
	for _, user := range users {
		userMap[user.Id] = user
	}

	var members []*dto.TeamMember
	for _, teamMember := range teamMembers {
		_, exists := userMap[teamMember.UserID]
		if !exists {
			continue
		}

		members = append(members, &dto.TeamMember{
			UserID: teamMember.UserID,
			TeamID: teamMember.TeamID,
			Role:   teamMember.Role,
		})
	}

	return members, nil
}
