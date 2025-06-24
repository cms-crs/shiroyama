package team

import "context"

func (service *Service) DeleteTeam(
	ctx context.Context,
	ID string,
) error {
	err := service.teamRepository.DeleteTeam(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}
