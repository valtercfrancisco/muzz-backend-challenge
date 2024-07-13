// Package repository provides implementations of repository interfaces for data access.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	explore "muzz-backend-challenge/pkg/proto"
)

// ExploreRepository defines methods for accessing exploration-related data.
type ExploreRepository interface {
	GetLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error)
	GetNewLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error)
	CountLikes(recipientUserID string) (int64, error)
	InsertDecision(actorUserID, recipientUserID string, likedRecipient bool) error
	InsertLike(actorUserID, recipientUserID string) error
	DeleteLike(actorUserID, recipientUserID string) error
	CheckMutualLike(actorUserID, recipientUserID string) (bool, error)
}

// exploreRepository implements the ExploreRepository interface.
type exploreRepository struct {
	db *sql.DB
}

// NewExploreRepository creates a new instance of exploreRepository.
func NewExploreRepository(db *sql.DB) ExploreRepository {
	return &exploreRepository{db: db}
}

// GetLikedYou retrieves a list of users who liked the recipient user.
func (r *exploreRepository) GetLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error) {
	query := `
        SELECT actor_user_id, EXTRACT(EPOCH FROM created_at) AS unix_timestamp
        FROM likes
        WHERE recipient_user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, recipientUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likers []*explore.ListLikedYouResponse_Liker
	for rows.Next() {
		var liker explore.ListLikedYouResponse_Liker
		var unixTimestamp float64 // Use float64 to handle double precision value from database
		if err := rows.Scan(&liker.ActorId, &unixTimestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		liker.UnixTimestamp = uint64(unixTimestamp) // Convert float64 to uint64

		likers = append(likers, &liker)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likers, nil
}

// GetNewLikedYou retrieves a list of new users who liked the recipient user.
func (r *exploreRepository) GetNewLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error) {
	query := `
        SELECT actor_user_id, EXTRACT(EPOCH FROM created_at) AS unix_timestamp
        FROM likes
        WHERE recipient_user_id = $1 
          AND actor_user_id NOT IN (
              SELECT recipient_user_id 
              FROM likes 
              WHERE actor_user_id = $1
          )
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, recipientUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var likers []*explore.ListLikedYouResponse_Liker
	for rows.Next() {
		var liker explore.ListLikedYouResponse_Liker
		var unixTimestamp float64
		if err := rows.Scan(&liker.ActorId, &unixTimestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		liker.UnixTimestamp = uint64(unixTimestamp)

		likers = append(likers, &liker)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return likers, nil
}

// CountLikes counts the number of users who liked the recipient user.
func (r *exploreRepository) CountLikes(recipientUserID string) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM likes WHERE recipient_user_id = $1"
	err := r.db.QueryRow(query, recipientUserID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// InsertDecision records a user's decision (like/dislike) regarding another user.
func (r *exploreRepository) InsertDecision(actorUserID, recipientUserID string, likedRecipient bool) error {
	query := `
        INSERT INTO decisions (actor_user_id, recipient_user_id, liked_recipient)
        VALUES ($1, $2, $3)
        ON CONFLICT (actor_user_id, recipient_user_id)
        DO UPDATE SET liked_recipient = EXCLUDED.liked_recipient, created_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, actorUserID, recipientUserID, likedRecipient)
	if err != nil {
		return fmt.Errorf("failed to insert decision: %w", err)
	}
	return nil
}

// InsertLike records a like action from the actor user to the recipient user.
func (r *exploreRepository) InsertLike(actorUserID, recipientUserID string) error {
	query := "INSERT INTO likes (actor_user_id, recipient_user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	_, err := r.db.Exec(query, actorUserID, recipientUserID)
	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}
	return nil
}

// DeleteLike removes a like action from the actor user to the recipient user.
func (r *exploreRepository) DeleteLike(actorUserID, recipientUserID string) error {
	query := "DELETE FROM likes WHERE actor_user_id = $1 AND recipient_user_id = $2"
	_, err := r.db.Exec(query, actorUserID, recipientUserID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}
	return nil
}

// CheckMutualLike checks if there is a mutual like between the actor user and the recipient user.
func (r *exploreRepository) CheckMutualLike(actorUserID, recipientUserID string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM likes WHERE actor_user_id = $1 AND recipient_user_id = $2)"
	var exists bool
	err := r.db.QueryRow(query, recipientUserID, actorUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check mutual like: %w", err)
	}
	return exists, nil
}
