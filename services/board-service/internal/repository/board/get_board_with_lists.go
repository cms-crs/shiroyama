package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
)

func (r *Repository) GetBoardWithLists(ctx context.Context, boardID string) (*boardv1.BoardWithLists, error) {
	board, err := r.GetBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	//lists, err := r.GetBoardLists(ctx, boardID)
	//if err != nil {
	//	return nil, err
	//}

	return &boardv1.BoardWithLists{
		Board: board,
		//Lists: lists,
	}, nil
}
