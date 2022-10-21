ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

generate:
	protoc -I . -I ${GOPATH}/src -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate --validate_out=lang=go,paths=source_relative:. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/users.proto

fmt:
	go fmt ./server

run:
	go run main.go