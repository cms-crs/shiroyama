package user

import (
	"context"
)

func (service *Service) DeleteUser(ctx context.Context, ID string) error {
	const op = "Service.DeleteUser"

	err := service.userRepository.DeleteUser(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}
