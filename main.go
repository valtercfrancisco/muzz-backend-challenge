package main

import (
	"context"
	"log"
	"muzz-backend-challenge/explore"
	"net"

	"google.golang.org/grpc"
)

type ExploreServiceImplementation struct {
	explore.UnimplementedExploreServiceServer
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

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal("cannot create lister: %s", err)
	}

	serviceRegistrar := grpc.NewServer()
	service := &ExploreServiceImplementation{}

	explore.RegisterExploreServiceServer(serviceRegistrar, service)
	err = serviceRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("Impossible to serve: %s", err)
	}
}
