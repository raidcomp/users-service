generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/users.proto

fmt:
	go fmt ./server

run:
	go run ./server/main.go