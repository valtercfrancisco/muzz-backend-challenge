generate_grpc_code:
	protoc --go_out=explore --go_opt=paths=source_relative --go-grpc_out=explore --go-grpc_opt=paths=source_relative explore-service.proto