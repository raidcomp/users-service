# users-service

Go-Lang gRPC microservice storing and managing users on [raidcomp.io](raidcomp.io).

## Technologies

- GoLang: microservice code
- gRPC/protobuf: service interface definition, protobuf client generation
- GitHub Actions: CI/CD
- Terraform: infrastructure as code
- AWS: cloud hosting provider

## Development

### `make run`

Starts the server for local development on `localhost:5785`

### `make generate`

Generates the gRPC server code based on the [users.proto](/proto/users.proto) definition.

### `make fmt`

Formats the go code (TODO: format terraform)

## Terraform

All AWS infrastructure is maintained in [/terraform](/terraform) directory. All terraform commands are run from here and require AWS account permissions to perform.