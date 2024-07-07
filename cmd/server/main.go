package main

import (
	"log"
	"muzz-backend-challenge/internal/config"
	"muzz-backend-challenge/internal/db"
	explore "muzz-backend-challenge/pkg/proto"
	"net"

	"muzz-backend-challenge/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := db.RunMigrations(dbConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal("cannot create listener: %s", err)
	}

	serviceRegistrar := grpc.NewServer()
	exploreService := service.NewExploreService()

	explore.RegisterExploreServiceServer(serviceRegistrar, exploreService)
	err = serviceRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("Impossible to serve: %s", err)
	}
}
