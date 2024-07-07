package service

import (
	"context"
	"fmt"
	explore "muzz-backend-challenge/pkg/proto"
	"muzz-backend-challenge/pkg/repository"
	"strconv"
)

type ExploreService struct {
	repository repository.ExploreRepository
	explore.UnimplementedExploreServiceServer
}

func NewExploreService(repo repository.ExploreRepository) *ExploreService {
	return &ExploreService{repository: repo}
}

func (service ExploreService) ListLikedYou(
	ctx context.Context,
	request *explore.ListLikedYouRequest,
) (*explore.ListLikedYouResponse, error) {
	recipientID := request.GetRecipientUserId()

	limit := 10
	offset := 0

	if request.GetPaginationToken() != "" {
		// Parse the pagination token if provided
		var err error
		offset, err = strconv.Atoi(request.GetPaginationToken())
		if err != nil {
			return nil, fmt.Errorf("invalid pagination token: %v", err)
		}
	}

	likers, err := service.repository.GetLikedYou(ctx, recipientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list likers: %v", err)
	}

	nextPaginationToken := fmt.Sprintf("%d", offset+limit)

	response := &explore.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: &nextPaginationToken,
	}

	return response, nil
}

func (service ExploreService) ListNewLikedYou(
	context.Context,
	*explore.ListLikedYouRequest,
) (*explore.ListLikedYouResponse, error) {
	return &explore.ListLikedYouResponse{}, nil
}

func (service ExploreService) CountLikedYou(
	context.Context,
	*explore.CountLikedYouRequest,
) (*explore.CountLikedYouResponse, error) {
	return &explore.CountLikedYouResponse{}, nil
}

func (service ExploreService) PutDecision(
	context.Context,
	*explore.PutDecisionRequest,
) (*explore.PutDecisionResponse, error) {
	return &explore.PutDecisionResponse{}, nil
}
