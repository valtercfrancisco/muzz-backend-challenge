// explore_service_test.go
package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	explore "muzz-backend-challenge/pkg/proto"
	"muzz-backend-challenge/pkg/repository"
)

func TestListLikedYou(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	recipientUserID := "test-recipient"
	paginationToken := "0"

	request := &explore.ListLikedYouRequest{
		RecipientUserId: recipientUserID,
		PaginationToken: &paginationToken,
	}

	likers := []*explore.ListLikedYouResponse_Liker{{ActorId: "user1"}}
	repo.On("GetLikedYou", ctx, recipientUserID, 10, 0).Return(likers, nil)

	response, err := service.ListLikedYou(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, likers, response.Likers)
	assert.Equal(t, "10", *response.NextPaginationToken)
	repo.AssertExpectations(t)
}

func TestListLikedYou_InvalidRecipientID(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	request := &explore.ListLikedYouRequest{}

	response, err := service.ListLikedYou(ctx, request)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "recipient user ID is required", status.Convert(err).Message())
}

func TestListLikedYou_InvalidPaginationToken(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	recipientID := "test-recipient"
	paginationToken := "invalid"

	request := &explore.ListLikedYouRequest{
		RecipientUserId: recipientID,
		PaginationToken: &paginationToken,
	}

	response, err := service.ListLikedYou(ctx, request)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "invalid pagination token")
}

func TestListNewLikedYou(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	recipientID := "test-recipient"
	paginationToken := "0"

	request := &explore.ListLikedYouRequest{
		RecipientUserId: recipientID,
		PaginationToken: &paginationToken,
	}

	likers := []*explore.ListLikedYouResponse_Liker{{ActorId: "user1"}}
	repo.On("GetNewLikedYou", ctx, recipientID, 10, 0).Return(likers, nil)

	response, err := service.ListNewLikedYou(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, likers, response.Likers)
	assert.Equal(t, "10", *response.NextPaginationToken)
	repo.AssertExpectations(t)
}

func TestListNewLikedYou_InvalidRecipientID(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	request := &explore.ListLikedYouRequest{}

	response, err := service.ListNewLikedYou(ctx, request)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "recipient user ID is required", status.Convert(err).Message())
}

func TestListNewLikedYou_InvalidPaginationToken(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	recipientID := "test-recipient"
	paginationToken := "invalid"

	request := &explore.ListLikedYouRequest{
		RecipientUserId: recipientID,
		PaginationToken: &paginationToken,
	}

	response, err := service.ListNewLikedYou(ctx, request)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "invalid pagination token")
}

func TestCountLikedYou(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	recipientID := "test-recipient"

	request := &explore.CountLikedYouRequest{
		RecipientUserId: recipientID,
	}

	count := 5
	repo.On("CountLikes", recipientID).Return(count, nil)

	response, err := service.CountLikedYou(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, uint64(count), response.Count)
	repo.AssertExpectations(t)
}

func TestCountLikedYou_Error(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	recipientID := "test-recipient"

	request := &explore.CountLikedYouRequest{
		RecipientUserId: recipientID,
	}

	repo.On("CountLikes", recipientID).Return(0, status.Errorf(codes.Internal, "count error"))

	response, err := service.CountLikedYou(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestPutDecision_LikedRecipient(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	repo.On("InsertDecision", actorID, recipientID, likedRecipient).Return(nil)
	repo.On("InsertLike", actorID, recipientID).Return(nil)
	repo.On("CheckMutualLike", actorID, recipientID).Return(true, nil)

	response, err := service.PutDecision(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.MutualLikes)
	repo.AssertExpectations(t)
}

func TestPutDecision_NotLikedRecipient(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := false

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	repo.On("InsertDecision", actorID, recipientID, likedRecipient).Return(nil)
	repo.On("DeleteLike", actorID, recipientID).Return(nil)

	response, err := service.PutDecision(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.MutualLikes)
	repo.AssertExpectations(t)
}

func TestPutDecision_InsertDecisionError(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	repo.On("InsertDecision", actorID, recipientID, likedRecipient).Return(status.Errorf(codes.Internal, "insert decision error"))

	response, err := service.PutDecision(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestPutDecision_InsertLikeError(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	repo.On("InsertDecision", actorID, recipientID, likedRecipient).Return(nil)
	repo.On("InsertLike", actorID, recipientID).Return(status.Errorf(codes.Internal, "insert like error"))

	response, err := service.PutDecision(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestPutDecision_CheckMutualLikeError(t *testing.T) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)

	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	repo.On("InsertDecision", actorID, recipientID, likedRecipient).Return(nil)
	repo.On("InsertLike", actorID, recipientID).Return(nil)
	repo.On("CheckMutualLike", actorID, recipientID).Return(false, status.Errorf(codes.Internal, "check mutual like error"))

	response, err := service.PutDecision(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
}
