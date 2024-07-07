package service

import (
	"context"
	explore "muzz-backend-challenge/pkg/proto"
)

type ExploreServiceImplementation struct {
	explore.UnimplementedExploreServiceServer
}

func NewExploreService() explore.ExploreServiceServer {
	return &ExploreServiceImplementation{}
}

func (s ExploreServiceImplementation) ListLikedYou(
	context.Context,
	*explore.ListLikedYouRequest,
) (*explore.ListLikedYouResponse, error) {
	return &explore.ListLikedYouResponse{
		Likers:              []*explore.ListLikedYouResponse_Liker{},
		NextPaginationToken: new(string),
	}, nil
}

func (s ExploreServiceImplementation) ListNewLikedYou(
	context.Context,
	*explore.ListLikedYouRequest,
) (*explore.ListLikedYouResponse, error) {
	return &explore.ListLikedYouResponse{}, nil
}

func (s ExploreServiceImplementation) CountLikedYou(
	context.Context,
	*explore.CountLikedYouRequest,
) (*explore.CountLikedYouResponse, error) {
	return &explore.CountLikedYouResponse{}, nil
}

func (s ExploreServiceImplementation) PutDecision(
	context.Context,
	*explore.PutDecisionRequest,
) (*explore.PutDecisionResponse, error) {
	return &explore.PutDecisionResponse{}, nil
}
