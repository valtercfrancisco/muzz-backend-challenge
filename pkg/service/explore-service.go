// Package service implements gRPC services for handling user exploration functionalities.
package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	explore "muzz-backend-challenge/pkg/proto"
	"muzz-backend-challenge/pkg/repository"
	"strconv"
)

// ExploreService implements the ExploreServiceServer interface.
type ExploreService struct {
	repository repository.ExploreRepository
	explore.UnimplementedExploreServiceServer
}

// NewExploreService creates a new instance of ExploreService.
func NewExploreService(repo repository.ExploreRepository) *ExploreService {
	return &ExploreService{repository: repo}
}

// ListLikedYou retrieves a list of users who liked the recipient user.
//
// It retrieves likers for the given recipient user ID. Pagination is supported
// via the pagination token, which allows fetching the next set of results.
func (service ExploreService) ListLikedYou(
	ctx context.Context,
	request *explore.ListLikedYouRequest,
) (*explore.ListLikedYouResponse, error) {
	recipientID := request.GetRecipientUserId()
	if recipientID == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient user ID is required")
	}

	limit := 10
	offset := 0

	if request.GetPaginationToken() != "" {
		// Parse the pagination token if provided
		var err error
		offset, err = strconv.Atoi(request.GetPaginationToken())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid pagination token: %v", err)
		}
	}

	likers, err := service.repository.GetLikedYou(ctx, recipientID, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list likers: %v", err)
	}

	nextPaginationToken := fmt.Sprintf("%d", offset+limit)

	response := &explore.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: &nextPaginationToken,
	}

	return response, nil
}

// ListNewLikedYou retrieves a list of new users who liked the recipient user. These are users who the recipient has seen or liked yet.
//
// It retrieves new likers for the given recipient user ID. Pagination is supported
// via the pagination token, which allows fetching the next set of results.
func (service ExploreService) ListNewLikedYou(
	ctx context.Context,
	request *explore.ListLikedYouRequest,
) (*explore.ListLikedYouResponse, error) {
	recipientID := request.GetRecipientUserId()
	if recipientID == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient user ID is required")
	}

	limit := 10
	offset := 0

	if request.GetPaginationToken() != "" {
		var err error
		offset, err = strconv.Atoi(request.GetPaginationToken())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid pagination token: %v", err)
		}
	}

	likers, err := service.repository.GetNewLikedYou(ctx, recipientID, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list likers: %v", err)
	}

	nextPaginationToken := fmt.Sprintf("%d", offset+limit)

	response := &explore.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: &nextPaginationToken,
	}

	return response, nil
}

// CountLikedYou counts the number of users who liked the recipient user.
//
// It returns the count of users who liked the recipient user specified in the request.
func (service ExploreService) CountLikedYou(
	_ context.Context,
	request *explore.CountLikedYouRequest,
) (*explore.CountLikedYouResponse, error) {
	count, err := service.repository.CountLikes(request.RecipientUserId)
	if err != nil {
		return nil, err
	}
	return &explore.CountLikedYouResponse{Count: uint64(count)}, nil
}

// PutDecision records a user's decision (like/dislike) regarding another user.
//
// It records the decision made by the actor user regarding the recipient user.
// If the decision results in a mutual like, it returns true in MutualLikes field.
func (service ExploreService) PutDecision(
	_ context.Context,
	request *explore.PutDecisionRequest,
) (*explore.PutDecisionResponse, error) {
	// Insert the decision into the decision database
	err := service.repository.InsertDecision(request.ActorUserId, request.RecipientUserId, request.LikedRecipient)
	if err != nil {
		return nil, err
	}

	mutualLikes := false

	// If the user liked the recipient, record the like
	if request.LikedRecipient {
		// Insert the like into the like database
		err = service.repository.InsertLike(request.ActorUserId, request.RecipientUserId)
		if err != nil {
			return nil, err
		}

		// Check if the recipient also liked the actor
		mutualLikes, err = service.repository.CheckMutualLike(request.ActorUserId, request.RecipientUserId)
		if err != nil {
			return nil, err
		}
	} else {
		// Delete the like if the actor passes on the recipient (unmatched)
		err = service.repository.DeleteLike(request.ActorUserId, request.RecipientUserId)
		if err != nil {
			return nil, err
		}
	}

	return &explore.PutDecisionResponse{MutualLikes: mutualLikes}, nil
}
