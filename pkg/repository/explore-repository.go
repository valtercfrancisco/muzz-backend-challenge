package repository

import (
	"context"
	"database/sql"
	"fmt"
	explore "muzz-backend-challenge/pkg/proto"
)

type ExploreRepository interface {
	GetLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error)
}

type exploreRepository struct {
	db *sql.DB
}

func NewExploreRepository(db *sql.DB) ExploreRepository {
	return &exploreRepository{db: db}
}

func (r *exploreRepository) GetLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error) {
	query := `
        SELECT actor_user_id, EXTRACT(EPOCH FROM created_at) AS unix_timestamp
        FROM likes
        WHERE recipient_user_id = $1 AND liked_recipient = TRUE
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
