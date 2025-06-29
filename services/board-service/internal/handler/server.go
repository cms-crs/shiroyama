package handler

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type BoardService interface {
	CreateBoard(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.Board, error)
	GetBoard(ctx context.Context, boardID string) (*boardv1.Board, error)
	GetBoardWithLists(ctx context.Context, boardID string) (*boardv1.BoardWithLists, error)
	UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.Board, error)
	DeleteBoard(ctx context.Context, boardID string) error
	GetUserBoards(ctx context.Context, userID string) ([]*boardv1.Board, error)
	GetTeamBoards(ctx context.Context, teamID string) ([]*boardv1.Board, error)
}

type ListService interface {
	CreateList(ctx context.Context, req *boardv1.CreateListRequest) (*boardv1.List, error)
	UpdateList(ctx context.Context, req *boardv1.UpdateListRequest) (*boardv1.List, error)
	GetList(ctx context.Context, listID string) (*boardv1.List, error)
	DeleteList(ctx context.Context, listID string) error
	ReorderLists(ctx context.Context, req *boardv1.ReorderListsRequest) error
}

type Handler struct {
	boardv1.UnimplementedBoardServiceServer
	log          *slog.Logger
	boardService BoardService
	listService  ListService
}

func NewHandler(log *slog.Logger, boardService BoardService, listService ListService) *Handler {
	return &Handler{
		log:          log,
		boardService: boardService,
		listService:  listService,
	}
}

func (h *Handler) CreateBoard(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.Board, error) {
	const op = "handler.CreateBoard"

	log := h.log.With(
		slog.String("op", op),
		slog.String("name", req.Name),
		slog.String("team_id", req.TeamId),
		slog.String("created_by", req.CreatedBy),
	)

	log.Info("Create board request received")

	board, err := h.boardService.CreateBoard(ctx, req)
	if err != nil {
		log.Error("Failed to create board", "error", err)
		return nil, err
	}

	log.Info("Board created successfully", "board_id", board.Id)

	return board, nil
}

func (h *Handler) GetBoard(ctx context.Context, req *boardv1.GetBoardRequest) (*boardv1.BoardWithLists, error) {
	const op = "handler.GetBoard"

	log := h.log.With(
		slog.String("op", op),
		slog.String("board_id", req.Id),
	)

	log.Info("Get board request received")

	board, err := h.boardService.GetBoardWithLists(ctx, req.Id)
	if err != nil {
		log.Error("Failed to get board", "error", err)
		return nil, err
	}

	log.Info("Board retrieved successfully")

	return board, nil
}

func (h *Handler) UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.Board, error) {
	const op = "handler.UpdateBoard"

	log := h.log.With(
		slog.String("op", op),
		slog.String("board_id", req.Id),
		slog.String("name", req.Name),
	)

	log.Info("Update board request received")

	board, err := h.boardService.UpdateBoard(ctx, req)
	if err != nil {
		log.Error("Failed to update board", "error", err)
		return nil, err
	}

	log.Info("Board updated successfully")

	return board, nil
}

func (h *Handler) DeleteBoard(ctx context.Context, req *boardv1.DeleteBoardRequest) (*emptypb.Empty, error) {
	const op = "handler.DeleteBoard"

	log := h.log.With(
		slog.String("op", op),
		slog.String("board_id", req.Id),
	)

	log.Info("Delete board request received")

	err := h.boardService.DeleteBoard(ctx, req.Id)
	if err != nil {
		log.Error("Failed to delete board", "error", err)
		return nil, err
	}

	log.Info("Board deleted successfully")

	return &emptypb.Empty{}, nil
}

func (h *Handler) GetUserBoards(ctx context.Context, req *boardv1.GetUserBoardsRequest) (*boardv1.GetUserBoardsResponse, error) {
	const op = "handler.GetUserBoards"

	log := h.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserId),
	)

	log.Info("Get user boards request received")

	boards, err := h.boardService.GetUserBoards(ctx, req.UserId)
	if err != nil {
		log.Error("Failed to get user boards", "error", err)
		return nil, err
	}

	log.Info("User boards retrieved successfully", "count", len(boards))

	return &boardv1.GetUserBoardsResponse{
		Boards: boards,
	}, nil
}

func (h *Handler) GetTeamBoards(ctx context.Context, req *boardv1.GetTeamBoardsRequest) (*boardv1.GetTeamBoardsResponse, error) {
	const op = "handler.GetTeamBoards"

	log := h.log.With(
		slog.String("op", op),
		slog.String("team_id", req.TeamId),
	)

	log.Info("Get team boards request received")

	boards, err := h.boardService.GetTeamBoards(ctx, req.TeamId)
	if err != nil {
		log.Error("Failed to get team boards", "error", err)
		return nil, err
	}

	log.Info("Team boards retrieved successfully", "count", len(boards))

	return &boardv1.GetTeamBoardsResponse{
		Boards: boards,
	}, nil
}

func (h *Handler) CreateList(ctx context.Context, req *boardv1.CreateListRequest) (*boardv1.List, error) {
	const op = "handler.CreateList"

	log := h.log.With(
		slog.String("op", op),
		slog.String("board_id", req.BoardId),
		slog.String("name", req.Name),
	)

	log.Info("Create list request received")

	list, err := h.listService.CreateList(ctx, req)
	if err != nil {
		log.Error("Failed to create list", "error", err)
		return nil, err
	}

	log.Info("List created successfully", "list_id", list.Id)

	return list, nil
}

func (h *Handler) UpdateList(ctx context.Context, req *boardv1.UpdateListRequest) (*boardv1.List, error) {
	const op = "handler.UpdateList"

	log := h.log.With(
		slog.String("op", op),
		slog.String("list_id", req.Id),
		slog.String("name", req.Name),
	)

	log.Info("Update list request received")

	list, err := h.listService.UpdateList(ctx, req)
	if err != nil {
		log.Error("Failed to update list", "error", err)
		return nil, err
	}

	log.Info("List updated successfully")

	return list, nil
}

func (h *Handler) DeleteList(ctx context.Context, req *boardv1.DeleteListRequest) (*emptypb.Empty, error) {
	const op = "handler.DeleteList"

	log := h.log.With(
		slog.String("op", op),
		slog.String("list_id", req.Id),
	)

	log.Info("Delete list request received")

	err := h.listService.DeleteList(ctx, req.Id)
	if err != nil {
		log.Error("Failed to delete list", "error", err)
		return nil, err
	}

	log.Info("List deleted successfully")

	return &emptypb.Empty{}, nil
}

func (h *Handler) ReorderLists(ctx context.Context, req *boardv1.ReorderListsRequest) (*emptypb.Empty, error) {
	const op = "handler.ReorderLists"

	log := h.log.With(
		slog.String("op", op),
		slog.String("board_id", req.BoardId),
		slog.Int("positions_count", len(req.Positions)),
	)

	log.Info("Reorder lists request received")

	err := h.listService.ReorderLists(ctx, req)
	if err != nil {
		log.Error("Failed to reorder lists", "error", err)
		return nil, err
	}

	log.Info("Lists reordered successfully")

	return &emptypb.Empty{}, nil
}

func Register(gRPCServer interface{}, log *slog.Logger, boardService BoardService, listService ListService) {
	handler := NewHandler(log, boardService, listService)
	boardv1.RegisterBoardServiceServer(gRPCServer.(*grpc.Server), handler)
}
