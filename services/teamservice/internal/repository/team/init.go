package team

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log/slog"
	"strings"
	"taskservice/internal/kafka"
	"time"
)

type Repository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewTeamRepository(log *slog.Logger, db *sql.DB) *Repository {
	return &Repository{log: log, db: db}
}

func (r *Repository) DeleteUserFromAllTeams(ctx context.Context, userID string) (*kafka.TeamDeletionData, error) {
	const op = "TeamRepository.DeleteUserFromAllTeams"

	r.log.Info("Deleting user from all teams", "user_id", userID, "op", op)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	defer func() {
		if r := recover(); r != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			panic(r)
		}
	}()

	query := `
        SELECT team_id, role, created_at 
        FROM team_members 
        WHERE user_id = $1 
        ORDER BY created_at
        `

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: failed to query user team memberships: %w", op, err)
	}
	defer rows.Close()

	var memberships []kafka.TeamMembershipData
	for rows.Next() {
		var membership kafka.TeamMembershipData
		var joinedAt time.Time

		if err := rows.Scan(&membership.TeamID, &membership.Role, &joinedAt); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s: failed to scan team membership: %w", op, err)
		}

		membership.JoinedAt = joinedAt
		memberships = append(memberships, membership)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: error iterating team memberships: %w", op, err)
	}

	deletionData := &kafka.TeamDeletionData{
		Teams: memberships,
	}

	deleteQuery := `DELETE FROM team_members WHERE user_id = $1`
	result, err := tx.ExecContext(ctx, deleteQuery, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: failed to delete user from teams: %w", op, err)
	}

	rowsAffected, _ := result.RowsAffected()

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	r.log.Info("Successfully deleted user from all teams",
		"user_id", userID,
		"teams_count", len(deletionData.Teams),
		"rows_affected", rowsAffected,
		"op", op)

	return deletionData, nil
}

func (r *Repository) RestoreUserTeams(ctx context.Context, userID string, data *kafka.TeamDeletionData) error {
	const op = "TeamRepository.RestoreUserTeams"

	r.log.Info("Restoring user teams", "user_id", userID, "teams_count", len(data.Teams), "op", op)

	if len(data.Teams) == 0 {
		r.log.Info("No teams to restore for user", "user_id", userID, "op", op)
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var restoredCount int
	var skippedCount int

	checkTeamQuery := `SELECT EXISTS(SELECT 1 FROM teams WHERE id = $1)`
	insertMemberQuery := `
        INSERT INTO team_members (user_id, team_id, role, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id, team_id) DO NOTHING`

	for _, team := range data.Teams {
		var teamExists bool
		if err := tx.QueryRowContext(ctx, checkTeamQuery, team.TeamID).Scan(&teamExists); err != nil {
			r.log.Warn("Failed to check team existence during restore",
				"team_id", team.TeamID,
				"user_id", userID,
				"error", err,
				"op", op)
			skippedCount++
			continue
		}

		if !teamExists {
			r.log.Warn("Team no longer exists, skipping restore",
				"team_id", team.TeamID,
				"user_id", userID,
				"op", op)
			skippedCount++
			continue
		}

		result, err := tx.ExecContext(ctx, insertMemberQuery,
			userID, team.TeamID, team.Role, team.JoinedAt, time.Now())
		if err != nil {
			r.log.Error("Failed to restore team membership",
				"team_id", team.TeamID,
				"user_id", userID,
				"error", err,
				"op", op)
			skippedCount++
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			restoredCount++
		} else {
			r.log.Debug("Team membership already exists, skipping",
				"team_id", team.TeamID,
				"user_id", userID,
				"op", op)
			skippedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	r.log.Info("Successfully restored user teams",
		"user_id", userID,
		"total_teams", len(data.Teams),
		"restored_count", restoredCount,
		"skipped_count", skippedCount,
		"op", op)

	return nil
}

func (r *Repository) GetUserTeamMemberships(ctx context.Context, userID string) (*kafka.TeamDeletionData, error) {
	const op = "TeamRepository.GetUserTeamMemberships"

	r.log.Info("Getting user team memberships", "user_id", userID, "op", op)

	query := `
        SELECT team_id, role, created_at 
        FROM team_members 
        WHERE user_id = $1 
        ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query user team memberships: %w", op, err)
	}
	defer rows.Close()

	var memberships []kafka.TeamMembershipData
	for rows.Next() {
		var membership kafka.TeamMembershipData
		var joinedAt time.Time

		if err := rows.Scan(&membership.TeamID, &membership.Role, &joinedAt); err != nil {
			return nil, fmt.Errorf("%s: failed to scan team membership: %w", op, err)
		}

		membership.JoinedAt = joinedAt
		memberships = append(memberships, membership)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating team memberships: %w", op, err)
	}

	deletionData := &kafka.TeamDeletionData{
		Teams: memberships,
	}

	r.log.Info("Retrieved user team memberships",
		"user_id", userID,
		"teams_count", len(deletionData.Teams),
		"op", op)

	return deletionData, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}

	errStr := strings.ToLower(err.Error())
	uniqueConstraintIndicators := []string{
		"duplicate key",
		"unique constraint",
		"unique_violation",
		"duplicate entry",
		"violates unique constraint",
		"already exists",
	}

	for _, indicator := range uniqueConstraintIndicators {
		if strings.Contains(errStr, indicator) {
			return true
		}
	}

	return false
}
