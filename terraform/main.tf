terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "us-east-1"
}

resource "aws_dynamodb_table" "users_dynamo_table" {
  name = "users"

  hash_key = "userID"

  write_capacity = 5
  read_capacity  = 5

  attribute {
    name = "userID"
    type = "S"
  }
}