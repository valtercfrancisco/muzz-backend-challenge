package repository

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/mock"
	explore "muzz-backend-challenge/pkg/proto"
)

type MockExploreRepository struct {
	mock.Mock
}

func (m *MockExploreRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockExploreRepository) GetLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error) {
	args := m.Called(ctx, recipientUserID, limit, offset)
	return args.Get(0).([]*explore.ListLikedYouResponse_Liker), args.Error(1)
}

func (m *MockExploreRepository) GetNewLikedYou(ctx context.Context, recipientUserID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error) {
	args := m.Called(ctx, recipientUserID, limit, offset)
	return args.Get(0).([]*explore.ListLikedYouResponse_Liker), args.Error(1)
}

func (m *MockExploreRepository) CountLikes(ctx context.Context, recipientUserID string) (int64, error) {
	args := m.Called(ctx, recipientUserID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockExploreRepository) InsertDecision(ctx context.Context, tx *sql.Tx, actorUserID, recipientUserID string, likedRecipient bool) error {
	args := m.Called(ctx, tx, actorUserID, recipientUserID, likedRecipient)
	return args.Error(0)
}

func (m *MockExploreRepository) InsertLike(ctx context.Context, tx *sql.Tx, actorUserID, recipientUserID string) error {
	args := m.Called(ctx, tx, actorUserID, recipientUserID)
	return args.Error(0)
}

func (m *MockExploreRepository) DeleteLike(ctx context.Context, tx *sql.Tx, actorUserID, recipientUserID string) error {
	args := m.Called(ctx, tx, actorUserID, recipientUserID)
	return args.Error(0)
}

func (m *MockExploreRepository) CheckMutualLike(ctx context.Context, tx *sql.Tx, actorUserID, recipientUserID string) (bool, error) {
	args := m.Called(ctx, tx, actorUserID, recipientUserID)
	return args.Bool(0), args.Error(1)
}
