package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
	explore "muzz-backend-challenge/pkg/proto"
)

type MockExploreRepository struct {
	mock.Mock
}

func (m *MockExploreRepository) GetLikedYou(ctx context.Context, recipientID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error) {
	args := m.Called(ctx, recipientID, limit, offset)
	return args.Get(0).([]*explore.ListLikedYouResponse_Liker), args.Error(1)
}

func (m *MockExploreRepository) GetNewLikedYou(ctx context.Context, recipientID string, limit, offset int) ([]*explore.ListLikedYouResponse_Liker, error) {
	args := m.Called(ctx, recipientID, limit, offset)
	return args.Get(0).([]*explore.ListLikedYouResponse_Liker), args.Error(1)
}

func (m *MockExploreRepository) CountLikes(recipientID string) (int64, error) {
	args := m.Called(recipientID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockExploreRepository) InsertDecision(actorID, recipientID string, liked bool) error {
	args := m.Called(actorID, recipientID, liked)
	return args.Error(0)
}

func (m *MockExploreRepository) InsertLike(actorID, recipientID string) error {
	args := m.Called(actorID, recipientID)
	return args.Error(0)
}

func (m *MockExploreRepository) CheckMutualLike(actorID, recipientID string) (bool, error) {
	args := m.Called(actorID, recipientID)
	return args.Bool(0), args.Error(1)
}

func (m *MockExploreRepository) DeleteLike(actorID, recipientID string) error {
	args := m.Called(actorID, recipientID)
	return args.Error(0)
}
