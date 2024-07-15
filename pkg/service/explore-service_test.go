package service

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	explore "muzz-backend-challenge/pkg/proto"
	"muzz-backend-challenge/pkg/repository"
)

func setupServiceAndRepo() (*repository.MockExploreRepository, *ExploreService) {
	repo := new(repository.MockExploreRepository)
	service := NewExploreService(repo)
	return repo, service
}

func TestListLikedYou(t *testing.T) {
	repo, service := setupServiceAndRepo()

	ctx := context.Background()
	recipientUserID := "test-recipient"
	paginationToken := "0"

	request := &explore.ListLikedYouRequest{
		RecipientUserId: recipientUserID,
		PaginationToken: &paginationToken,
	}

	likers := []*explore.ListLikedYouResponse_Liker{{ActorId: "user1"}}
	repo.On("GetLikedYou", mock.Anything, recipientUserID, 10, 0).Return(likers, nil)

	response, err := service.ListLikedYou(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, likers, response.Likers)
	assert.Equal(t, "10", *response.NextPaginationToken)
	repo.AssertExpectations(t)
}

func TestListLikedYou_InvalidRecipientID(t *testing.T) {
	_, service := setupServiceAndRepo()

	ctx := context.Background()
	request := &explore.ListLikedYouRequest{}

	response, err := service.ListLikedYou(ctx, request)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "recipient user ID is required", status.Convert(err).Message())
}

func TestListLikedYou_InvalidPaginationToken(t *testing.T) {
	_, service := setupServiceAndRepo()

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
	repo, service := setupServiceAndRepo()

	ctx := context.Background()
	recipientID := "test-recipient"
	paginationToken := "0"

	request := &explore.ListLikedYouRequest{
		RecipientUserId: recipientID,
		PaginationToken: &paginationToken,
	}

	likers := []*explore.ListLikedYouResponse_Liker{{ActorId: "user1"}}
	repo.On("GetNewLikedYou", mock.Anything, recipientID, 10, 0).Return(likers, nil)

	response, err := service.ListNewLikedYou(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, likers, response.Likers)
	assert.Equal(t, "10", *response.NextPaginationToken)
	repo.AssertExpectations(t)
}

func TestListNewLikedYou_InvalidRecipientID(t *testing.T) {
	_, service := setupServiceAndRepo()

	ctx := context.Background()
	request := &explore.ListLikedYouRequest{}

	response, err := service.ListNewLikedYou(ctx, request)

	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "recipient user ID is required", status.Convert(err).Message())
}

func TestListNewLikedYou_InvalidPaginationToken(t *testing.T) {
	_, service := setupServiceAndRepo()

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
	repo, service := setupServiceAndRepo()
	ctx := context.Background()
	recipientID := "recipient-user"
	expectedCount := int64(5)

	repo.On("CountLikes", ctx, recipientID).Return(expectedCount, nil)

	request := &explore.CountLikedYouRequest{RecipientUserId: recipientID}
	response, err := service.CountLikedYou(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, uint64(expectedCount), response.Count)
	repo.AssertExpectations(t)
}

func TestCountLikedYou_Error(t *testing.T) {
	repo, service := setupServiceAndRepo()

	ctx := context.Background()
	recipientID := "test-recipient"

	request := &explore.CountLikedYouRequest{
		RecipientUserId: recipientID,
	}

	repo.On("CountLikes", mock.Anything, recipientID).Return(int64(0), status.Errorf(codes.Internal, "count error"))

	response, err := service.CountLikedYou(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestPutDecision_LikedRecipient(t *testing.T) {
	repo, service := setupServiceAndRepo()
	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	// Use sqlmock to create a valid mock transaction
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mockTx, err := db.Begin()
	assert.NoError(t, err)

	// Mocking the repository methods
	repo.On("BeginTransaction", ctx).Return(mockTx, nil)
	repo.On("InsertDecision", ctx, mockTx, actorID, recipientID, likedRecipient).Return(nil)
	repo.On("InsertLike", ctx, mockTx, actorID, recipientID).Return(nil)
	repo.On("CheckMutualLike", ctx, mockTx, actorID, recipientID).Return(true, nil)

	// Expect the transaction to commit
	mock.ExpectCommit()

	response, err := service.PutDecision(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.MutualLikes)
	repo.AssertExpectations(t)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPutDecision_NotLikedRecipient(t *testing.T) {
	repo, service := setupServiceAndRepo()
	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := false

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	// Use sqlmock to create a valid mock transaction
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mockTx, err := db.Begin()
	assert.NoError(t, err)

	// Mocking the repository methods
	repo.On("BeginTransaction", ctx).Return(mockTx, nil)
	repo.On("InsertDecision", ctx, mockTx, actorID, recipientID, likedRecipient).Return(nil)
	repo.On("DeleteLike", ctx, mockTx, actorID, recipientID).Return(nil)

	// Expect the transaction to commit
	mock.ExpectCommit()

	response, err := service.PutDecision(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.MutualLikes)
	repo.AssertExpectations(t)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPutDecision_InsertDecisionError(t *testing.T) {
	repo, service := setupServiceAndRepo()
	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	// Use sqlmock to create a valid mock transaction
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mockTx, err := db.Begin()
	assert.NoError(t, err)

	// Mocking the repository methods
	repo.On("BeginTransaction", ctx).Return(mockTx, nil)
	repo.On("InsertDecision", ctx, mockTx, actorID, recipientID, likedRecipient).Return(status.Errorf(codes.Internal, "insert decision error"))

	// Expect the transaction to rollback
	mock.ExpectRollback()

	response, err := service.PutDecision(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
	repo.AssertExpectations(t)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPutDecision_InsertLikeError(t *testing.T) {
	repo, service := setupServiceAndRepo()
	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	// Use sqlmock to create a valid mock transaction
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mockTx, err := db.Begin()
	assert.NoError(t, err)

	// Mocking the repository methods
	repo.On("BeginTransaction", ctx).Return(mockTx, nil)
	repo.On("InsertDecision", ctx, mockTx, actorID, recipientID, likedRecipient).Return(nil)
	repo.On("InsertLike", ctx, mockTx, actorID, recipientID).Return(status.Errorf(codes.Internal, "insert like error"))

	// Expect the transaction to rollback
	mock.ExpectRollback()

	response, err := service.PutDecision(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
	repo.AssertExpectations(t)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPutDecision_CheckMutualLikeError(t *testing.T) {
	repo, service := setupServiceAndRepo()
	ctx := context.Background()
	actorID := "actor-user"
	recipientID := "recipient-user"
	likedRecipient := true

	request := &explore.PutDecisionRequest{
		ActorUserId:     actorID,
		RecipientUserId: recipientID,
		LikedRecipient:  likedRecipient,
	}

	// Use sqlmock to create a valid mock transaction
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mockTx, err := db.Begin()
	assert.NoError(t, err)

	// Mocking the repository methods
	repo.On("BeginTransaction", ctx).Return(mockTx, nil)
	repo.On("InsertDecision", ctx, mockTx, actorID, recipientID, likedRecipient).Return(nil)
	repo.On("InsertLike", ctx, mockTx, actorID, recipientID).Return(nil)
	repo.On("CheckMutualLike", ctx, mockTx, actorID, recipientID).Return(false, status.Errorf(codes.Internal, "check mutual like error"))

	// Expect the transaction to rollback
	mock.ExpectRollback()

	response, err := service.PutDecision(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.Internal, status.Code(err))
	repo.AssertExpectations(t)
	assert.NoError(t, mock.ExpectationsWereMet())
}
