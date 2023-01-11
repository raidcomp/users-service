# users-service

Go-Lang gRPC microservice storing and managing users on [raidcomp.io](raidcomp.io).

## Technologies

- GoLang: microservice code
- gRPC: service interface definition, protobuf client generation
- GitHub Actions: CI/CD
- AWS: cloud hosting provider

## Development

### `make run`

Starts the server for local development on `localhost:5785`

### `make generate`

Generates the gRPC server code based on the [users.proto](/proto/users.proto) definition.

### `make fmt`

Formats the go code (TODO: format terraform)
